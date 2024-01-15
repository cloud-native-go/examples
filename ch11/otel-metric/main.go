/*
 * Copyright 2024 Matthew A. Titmus
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
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

const (
	serviceName = "fibonacci"
)

// The requests counter instrument. As a synchronous instrument,
// we'll need to keep it so we can use it later to record data.
var requests metric.Int64Counter

// Define our labels here so that we can easily reuse them.
var attributes = []attribute.KeyValue{
	attribute.Key("application").String(serviceName),
	attribute.Key("container_id").String(os.Getenv("HOSTNAME")),
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	ctx := context.Background()

	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}

	// Now we can register it as the otel meter provider.
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(exporter))
	defer provider.Shutdown(ctx)

	meter := provider.Meter(serviceName)

	if err := buildRequestsCounter(meter); err != nil {
		log.Fatal(err)
	}
	if err := buildRuntimeObservers(meter); err != nil {
		log.Fatal(err)
	}

	log.Println("Browse to localhost:3000?n=6")

	// Neat, huh?
	http.Handle("/metrics", promhttp.Handler())
	http.Handle("/", http.HandlerFunc(fibHandler))

	log.Fatal(http.ListenAndServe(":3000", nil))
}

func fibHandler(w http.ResponseWriter, req *http.Request) {
	var err error
	var n int

	if len(req.URL.Query()["n"]) != 1 {
		err = fmt.Errorf("wrong number of arguments")
	} else {
		n, err = strconv.Atoi(req.URL.Query()["n"][0])
	}

	if err != nil {
		http.Error(w, "couldn't parse index n", http.StatusBadRequest)
		return
	}

	ctx := req.Context()

	// Call the child function, passing it the request context.
	result := <-Fibonacci(ctx, n)

	fmt.Fprintln(w, result)
}

func buildRequestsCounter(meter metric.Meter) error {
	var err error

	// Get an Int64Counter for a metric called "fibonacci_requests_total".
	requests, err = meter.Int64Counter("fibonacci_requests_total",
		metric.WithDescription("Total number of Fibonacci requests."),
		metric.WithUnit("{request}"),
	)

	return err
}

func buildRuntimeObservers(meter metric.Meter) error {
	var err error
	var m runtime.MemStats

	_, err = meter.Int64ObservableUpDownCounter("fibonacci_memory_usage_bytes",
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			runtime.ReadMemStats(&m)
			o.Observe(int64(m.Sys), metric.WithAttributes(attributes...))
			return nil
		}),
		metric.WithDescription("Amount of memory used."),
		metric.WithUnit("By"),
	)
	if err != nil {
		return err
	}

	_, err = meter.Int64ObservableGauge(
		"fibonacci_num_goroutines",
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			o.Observe(int64(runtime.NumGoroutine()), metric.WithAttributes(attributes...))
			return nil
		}),
		metric.WithDescription("Number of running goroutines."),
		metric.WithUnit("{goroutine}"),
	)
	if err != nil {
		return err
	}

	return nil
}

func Fibonacci(ctx context.Context, n int) chan int {
	requests.Add(ctx, 1, metric.WithAttributes(attributes...))

	ch := make(chan int)

	go func() {
		result := 1
		if n > 1 {
			a := Fibonacci(ctx, n-1)
			b := Fibonacci(ctx, n-2)
			result = <-a + <-b
		}

		ch <- result
	}()

	return ch
}
