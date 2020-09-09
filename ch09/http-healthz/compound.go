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
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

type Health struct {
	Database     bool
	IndexService bool
}

func healthCompoundHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the context from the request and add a 5-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	health := &Health{}    // Create our health check data wrapper
	wg := sync.WaitGroup{} // A WaitGroup helps us check concurrently

	// Probe #1: A database functionality check
	wg.Add(1)
	go func() {
		defer wg.Done()
		probeDatabase(ctx, health)
	}()

	// Probe #2: A downstream service functionality check
	wg.Add(1)
	go func() {
		defer wg.Done()
		probeIndexService(ctx, health)
	}()

	wg.Wait()

	// Marshal our health struct into JSON and return it
	bytes, err := json.Marshal(health)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(bytes)
}

func probeDatabase(ctx context.Context, health *Health) {
	// An imaginary function that executes a simple database query.
	if err := HealthCheck(ctx); err == nil {
		health.Database = true
	} else {
		log.Println(err)
	}
}

func probeIndexService(ctx context.Context, health *Health) {
	// An imaginary function that executes a simple database query.
	if err := HealthCheck(ctx); err == nil {
		health.IndexService = true
	} else {
		log.Println(err)
	}
}
