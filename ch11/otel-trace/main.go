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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	jaegerEndpoint = "localhost:4317"
	serviceName    = "Fibonacci"
)

var tracer trace.Tracer

func createAndRegisterExporters(ctx context.Context) error {
	// Create and configure the stdout exporter
	stdExporter, err := stdouttrace.New(
		stdouttrace.WithPrettyPrint(),
	)
	if err != nil {
		return err
	}

	// Create and configure the Jaeger exporter
	otlpExporter, err := otlptracegrpc.New(
		ctx,
		otlptracegrpc.WithEndpoint(jaegerEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return err
	}

	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return err
	}

	// Create and configure the TracerProvider exporter using the
	// newly-created exporters.
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(stdExporter),
		sdktrace.WithBatcher(otlpExporter),
		sdktrace.WithResource(r),
	)

	// Now we can register tp as the otel trace provider.
	// otel.SetTracerProvider(tp)

	// Finally, set the tracer that can be used for this package.
	tracer = tp.Tracer(serviceName)

	return nil
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	err := createAndRegisterExporters(ctx)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}

	fmt.Println("Browse to localhost:3000?n=6")

	http.Handle("/", otelhttp.NewHandler(http.HandlerFunc(fibHandler), "root"))

	if err := http.ListenAndServe(":3000", nil); err != nil {
		slog.ErrorContext(ctx, err.Error())
		return
	}
}

func fibHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var n int

	ctx := req.Context()

	// Get the Span associated with the current context and
	// attach the parameter and result as attributes.
	sp := trace.SpanFromContext(ctx)

	if len(req.URL.Query()["n"]) != 1 {
		err = fmt.Errorf("wrong number of arguments")
	} else {
		n, err = strconv.Atoi(req.URL.Query()["n"][0])
	}

	if err != nil {
		http.Error(w, "couldn't parse index n", http.StatusBadRequest)
		return
	}

	sp.SetAttributes(attribute.Int("parameter", n))

	// Call the child function, passing it the request context.
	result := <-Fibonacci(ctx, n)

	sp.SetAttributes(attribute.Int("result", result))

	fmt.Fprintln(w, result)
}

func Fibonacci(ctx context.Context, n int) chan int {
	ch := make(chan int)

	go func() {
		ctx, sp := tracer.Start(ctx,
			"fibonacci",
			trace.WithAttributes(
				attribute.Int("n", n)),
		)
		defer sp.End()

		result := 1
		if n > 1 {
			a := Fibonacci(ctx, n-1)
			b := Fibonacci(ctx, n-2)
			result = <-a + <-b
		}

		sp.SetAttributes(attribute.Int("result", result))

		ch <- result
	}()

	return ch
}
