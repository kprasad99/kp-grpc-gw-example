package io.github.kprasad99.example.interceptor;

import java.nio.charset.StandardCharsets;
import java.util.Optional;

import org.apache.logging.log4j.util.Strings;
import org.lognet.springboot.grpc.GRpcGlobalInterceptor;
import org.springframework.core.Ordered;

import io.grpc.Metadata;
import io.grpc.ServerCall;
import io.grpc.ServerCallHandler;
import io.grpc.ServerInterceptor;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.api.trace.Span;
import io.opentelemetry.api.trace.SpanBuilder;
import io.opentelemetry.api.trace.SpanContext;
import io.opentelemetry.api.trace.TraceFlags;
import io.opentelemetry.api.trace.TraceState;
import io.opentelemetry.api.trace.Tracer;
import io.opentelemetry.context.Context;
import io.opentelemetry.context.Scope;
import lombok.AllArgsConstructor;
import lombok.RequiredArgsConstructor;
import lombok.Setter;
import lombok.experimental.Accessors;
import lombok.extern.slf4j.Slf4j;

@Slf4j
@RequiredArgsConstructor
public class TracingServerInterceptor implements ServerInterceptor, Ordered {

    @Setter
    @Accessors(fluent = true)
    private Integer order;

    private final Tracer tracer;

    @Override
    public <ReqT, RespT> ServerCall.Listener<ReqT> interceptCall(
            ServerCall<ReqT, RespT> call, Metadata headers, ServerCallHandler<ReqT, RespT> next) {

    	headers.keys().stream().forEach(log::info);
    	
        String traceId = headers.get(Metadata.Key.of("traceId", Metadata.ASCII_STRING_MARSHALLER));
        if(Strings.isEmpty(traceId)) {
        	traceId = headers.get(Metadata.Key.of("traceparent", Metadata.ASCII_STRING_MARSHALLER));
        	log.info("TraceID {}", traceId);
        }
        String spanId = headers.get(Metadata.Key.of("spanId", Metadata.ASCII_STRING_MARSHALLER));
        if(Strings.isEmpty(spanId)) {
        	spanId = headers.get(Metadata.Key.of("tracespan", Metadata.ASCII_STRING_MARSHALLER));
        }
        
        Context spanContext = createSpanContext(traceId, spanId);
        SpanBuilder spanBuilder = tracer.spanBuilder("grpc-spring-boot-starter-span").setParent(spanContext);
        Span span = spanBuilder.startSpan();
        Context scopedContext = Context.current().with(span);

        try (Scope scope = scopedContext.makeCurrent()) {
            return next.startCall(call, headers);
        } finally {
            span.end();
        }
    }

    private Context createSpanContext(String traceId, String spanId) {
        SpanContext spanContext = SpanContext.createFromRemoteParent(
                traceId != null ? traceId : "",
                spanId != null ? spanId : "",
                TraceFlags.getDefault(),
                TraceState.getDefault()
        );
        log.info("SpanContext trace={},span={}", spanContext.getTraceId(), spanContext.getSpanId());
        return Context.current().with(Span.wrap(spanContext));
    }

    @Override
    public int getOrder() {
        return Optional.ofNullable(order).orElse(HIGHEST_PRECEDENCE);
    }
}
