package io.github.kprasad99.example.configuration;

import org.lognet.springboot.grpc.GRpcGlobalInterceptor;
import org.springframework.boot.actuate.autoconfigure.tracing.ConditionalOnEnabledTracing;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

import io.grpc.ServerInterceptor;
import io.opentelemetry.api.OpenTelemetry;
import io.opentelemetry.instrumentation.grpc.v1_6.GrpcTelemetry;
import lombok.extern.slf4j.Slf4j;


@Slf4j
@Configuration
@ConditionalOnEnabledTracing
public class GrpcConfiguration {

//	@Bean
//    @GRpcGlobalInterceptor
//    public TracingServerInterceptor otelGrpcInterceptor(Tracer tracer) {
//		log.info("Enabling otel tracer for grpc");
//		return new TracingServerInterceptor(tracer);
//	}
	
	@Bean
    @GRpcGlobalInterceptor
    public ServerInterceptor otelGrpcInterceptor(OpenTelemetry telemetry) {
		log.info("Enabling otel tracer for grpc");
		return GrpcTelemetry.create(telemetry).newServerInterceptor();
	}
	
}
