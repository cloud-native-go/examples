/*
 * Copyright 2020 Matthew A. Titmus
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
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/metric/prometheus"
	"go.opentelemetry.io/otel/label"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/unit"
)

const (
	jaegerEndpoint = "http://localhost:14268/api/traces"
	serviceName    = "fibonacci"
)

// The requests counter instrument. As a synchronous instrument,
// we'll need to keep it so we can use it later to record data.
var requests metric.Int64Counter

// Define our labels here so that we can easily reuse them.
var labels = []label.KeyValue{
	label.Key("application").String(serviceName),
	label.Key("container_id").String(os.Getenv("HOSTNAME")),
}

func init() {
	log.SetFlags(log.Lshortfile)
}

func main() {
	// Create and configure the Prometheus exporter
	promExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Now we can register it as the otel meter provider.
	otel.SetMeterProvider(promExporter.MeterProvider())

	if err := buildRequestsCounter(); err != nil {
		log.Fatal(err)
	}
	if err := buildRuntimeObservers(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			log.Println(requestsCount)
			time.Sleep(time.Second)
		}
	}()

	log.Println("Browse to localhost:3000?n=6")

	// Neat, huh?
	http.Handle("/metrics", promExporter)
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
		http.Error(w, "couldn't parse index n", 400)
		return
	}

	ctx := req.Context()

	// Call the child function, passing it the request context.
	result := <-Fibonacci(ctx, n)

	fmt.Fprintln(w, result)
}

func parseArguments() (int, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return 0, fmt.Errorf("expected an int argument")
	}

	n, err := strconv.Atoi(args[0])
	if err != nil {
		return 0, fmt.Errorf("can't parse argument as integer: %w", err)
	}

	return n, nil
}

func buildRequestsCounter() error {
	var err error

	// Retrieve the meter from the meter provider.
	meter := otel.GetMeterProvider().Meter(serviceName)

	// Get an Int64Counter for a metric called "fibonacci_requests_total".
	requests, err = meter.NewInt64Counter("fibonacci_requests_total",
		metric.WithDescription("Total number of Fibonacci requests."),
	)

	return err
}

func buildRuntimeObservers() error {
	var err error
	var m runtime.MemStats

	meter := otel.GetMeterProvider().Meter(serviceName)

	_, err = meter.NewInt64UpDownSumObserver("memory_usage_bytes",
		func(_ context.Context, result metric.Int64ObserverResult) {
			runtime.ReadMemStats(&m)
			log.Println("memory_usage_bytes", int64(m.Sys))
			result.Observe(int64(m.Sys), labels...)
		},
		metric.WithDescription("Amount of memory used."),
		metric.WithUnit(unit.Bytes),
	)
	if err != nil {
		return err
	}

	_, err = meter.NewInt64UpDownSumObserver("num_goroutines",
		func(_ context.Context, result metric.Int64ObserverResult) {
			log.Println("num_goroutines", int64(runtime.NumGoroutine()))
			result.Observe(int64(runtime.NumGoroutine()), labels...)
		},
		metric.WithDescription("Number of running goroutines."),
	)
	if err != nil {
		return err
	}

	return nil
}

var requestsCount int64

func Fibonacci(ctx context.Context, n int) chan int {
	requests.Add(ctx, 1, labels...)

	atomic.AddInt64(&requestsCount, 1)

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
