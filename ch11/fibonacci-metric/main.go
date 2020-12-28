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

var requestsCount int64

func init() {
	log.SetFlags(log.Lshortfile)
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

func main() {
	// Create and configure the Prometheus exporter
	promExporter, err := prometheus.NewExportPipeline(prometheus.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// Now we can register it as the otel meter provider.
	otel.SetMeterProvider(promExporter.MeterProvider())

	go updateMetrics(context.Background())
	buildResultsCounterObserver()

	fmt.Println("Browse to localhost:3000?n=6")

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

func Fibonacci(ctx context.Context, n int) chan int {
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

var metricLabels = []label.KeyValue{
	label.Key("application").String(serviceName),
	label.Key("container_id").String(os.Getenv("HOSTNAME")),
}

func buildResultsCounterObserver() {
	callbackFunc := func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(requestsCount, metricLabels...)
	}

	meter := otel.GetMeterProvider().Meter(serviceName)

	meter.NewInt64SumObserver("fibonacci_requests_total",
		callbackFunc,
		metric.WithDescription("Count of Fibonacci requests."),
	)
}

func updateMetrics(ctx context.Context) {
	meter := otel.GetMeterProvider().Meter(serviceName)

	mem, _ := meter.NewInt64UpDownCounter("memory_usage_bytes",
		metric.WithDescription("Amount of memory used."),
		metric.WithUnit(unit.Bytes),
	)
	goroutines, _ := meter.NewInt64UpDownCounter("num_goroutines",
		metric.WithDescription("Number of running goroutines."),
	)

	var m runtime.MemStats

	for {
		runtime.ReadMemStats(&m)

		fmt.Println(m.Sys, runtime.NumGoroutine())

		mMem := mem.Measurement(int64(m.Sys))
		mGoroutines := goroutines.Measurement(int64(runtime.NumGoroutine()))

		meter.RecordBatch(ctx, metricLabels, mMem, mGoroutines)

		time.Sleep(5 * time.Second)
	}
}
