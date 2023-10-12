package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gofiber/contrib/otelfiber/v2"
	"github.com/gofiber/fiber/v2"
	service "github.com/kprasad99/kp-grpc-gateway/api/rest"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func okHandler(ctx *fiber.Ctx) error {
	return ctx.SendString("OK")
}

var tracer = otel.Tracer("kp-grpc-gw")

func main() {
	contextPath := os.Getenv("WEB_CONTEXT_PATH")
	port := os.Getenv("WEB_PORT")
	host := os.Getenv("WEB_HOST")
	backendServer := os.Getenv("BACKEND_URL")
	otelAddress := os.Getenv("OTEL_ADDRESS")
	if port == "" {
		port = "8080"
	}

	address := ""
	if host == "" {
		address = ":" + port
	} else {
		address = host + ":" + port
	}

	contextPath = strings.TrimSuffix(contextPath, "/")

	app := fiber.New()
	app.Get("/liveness", okHandler)
	app.Get("/readiness", okHandler)

	var conn *grpc.ClientConn
	var err error
	if otelAddress == "" {
		conn, err = grpc.Dial(backendServer, grpc.WithTransportCredentials(insecure.NewCredentials()))
	} else {
		conn, err = grpc.Dial(backendServer,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
			grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()))
	}

	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	userService := service.NewUserService(conn)

	router := app.Group("/api")

	if otelAddress != "" {
		tp, err := initProvider(otelAddress)
		if err != nil {
			log.Panic("Failed to initialize otel", err)
		}
		defer func() {
			if err := tp.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}()
		app.Use(otelfiber.Middleware())

	}

	if contextPath != "" {
		router = app.Group(contextPath + "/api")
	}
	router.Get("/users", userService.List)
	router.Post("/users/:id", userService.Add)
	router.Put("/users/:id", userService.Update)
	router.Delete("/users/:id", userService.Remove)
	log.Println("Started web server")
	if err := app.Listen(address); err != nil {
		log.Panic("Failed to start web server", err)
	}

}

func initProvider(otelAddress string) (*trace.TracerProvider, error) {
	ctx := context.Background()

	res, err := resource.New(ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceName("kp-grpc-gw"),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// If the OpenTelemetry Collector is running on a local cluster (minikube or
	// microk8s), it should be accessible through the NodePort service at the
	// `localhost:30080` endpoint. Otherwise, replace `localhost` with the
	// endpoint of your cluster. If you run the app inside k8s, then you can
	// probably connect directly to the service through dns.
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, otelAddress,
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	// Set up a trace exporter
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}

	// Register the trace exporter with a TracerProvider, using a batch
	// span processor to aggregate spans before export.
	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tracerProvider)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})

	// Shutdown will flush any remaining spans and shut down the exporter.
	return tracerProvider, nil
}
