/*
 * Copyright 2023 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	serviceName    = "Fibonacci"
	serviceVersion = "0.0.2"

	jaegerEndpoint = "localhost:4317"
)

var tracer trace.Tracer

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(serviceName),
		semconv.ServiceVersion(serviceVersion),
	)

	// Ensure default SDK resources and the required service name are set
	res, err := resource.Merge(resource.Default(), res)
	if err != nil {
		return nil, fmt.Errorf("failed to merge resources: %w", err)
	}

	// Create and configure the stdout exporter
	stdExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build StdoutExporter: %w", err)
	}

	// Create and configure the OTLP exporter for Jaeger
	otlpExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build OtlpExporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(stdExporter),
		sdktrace.WithBatcher(otlpExporter),
	)

	// Create and configure the TracerProvider exporter using the
	// newly-created exporters.
	return tp, nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	tp, err := newTracerProvider(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	// Handle shutdown properly so nothing leaks.
	defer func() { tp.Shutdown(ctx) }()

	// Registers tp as the global trace provider to allow
	// auto-instrumentation to access it
	otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(serviceName)

	fmt.Println("Browse to localhost:3000?n=6")

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(fibHandler), "root"))

	if err := http.ListenAndServe(":3000", nil); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
}

func fibHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	// Get the Span associated with the current context and
	// attach the parameter and result as attributes.
	sp := trace.SpanFromContext(ctx)

	args := req.URL.Query()["n"]
	if len(args) != 1 {
		msg := "wrong number of arguments"
		sp.SetStatus(codes.Error, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	sp.SetAttributes(attribute.String("fibonacci.argument", args[0]))

	n, err := strconv.Atoi(args[0])
	if err != nil {
		msg := fmt.Sprintf("couldn't parse index n: %s", err.Error())
		sp.SetStatus(codes.Error, msg)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	sp.SetAttributes(attribute.Int("fibonacci.parameter", n))

	// Call the child function, passing it the request context.
	result := Fibonacci(ctx, n)

	sp.SetAttributes(attribute.Int("fibonacci.result", result))

	// Finally, send the result back in the response.
	fmt.Fprintln(w, result)
}

func Fibonacci(ctx context.Context, n int) int {
	ctx, sp := tracer.Start(ctx,
		"fibonacci",
		trace.WithAttributes(attribute.Int("fibonacci.n", n)),
	)
	defer sp.End()

	result := 1
	if n > 1 {
		a := Fibonacci(ctx, n-1)
		b := Fibonacci(ctx, n-2)
		result = a + b
	}

	sp.SetAttributes(attribute.Int("fibonacci.result", result))

	return result
}
