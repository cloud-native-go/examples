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
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Get a context that closes on SIGTERM, SIGINT, or SIGQUIT
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	defer cancel()

	server := &http.Server{Addr: ":8080"}

	// Register a cleanup function to be automatically
	// called when the server is shut down
	server.RegisterOnShutdown(doCleanup)

	// Register the readiness and liveness probes.
	http.Handle("/ready", handleReadiness(ctx))
	http.Handle("/health", handleLiveness())

	// This goroutine will respond to context closure
	// by shutting down the server
	go func() {
		// Read from the context's Done channel
		// This operation will block until the context closes
		<-ctx.Done()

		log.Println("Got shutdown signal.")

		// Wait for the readiness probe to detect the failure
		<-time.After(5 * time.Second)

		// Issue the shutdown proper. Don't pass the
		// already-closed Context value to it!
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Error while stopping HTTP listener: %s", err)
		}
	}()

	// Begin listening on :8080
	log.Println(server.ListenAndServe())
}

func doCleanup() {
	log.Println("Cleanup starting.")
	time.Sleep(5 * time.Second)
	log.Println("Cleanup complete.")
}

func handleLiveness() http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	return http.HandlerFunc(f)
}

func handleReadiness(ctx context.Context) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusServiceUnavailable)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}
	return http.HandlerFunc(f)
}
