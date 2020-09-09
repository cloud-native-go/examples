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
	"net/http"
	"time"
)

var service Service

func healthDeepHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the context from the request and add a 5-second timeout
	ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
	defer cancel()

	// Can the service execute a key query against the database?
	if err := service.GetUser(ctx); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type Service struct{}

func (s Service) GetUser(ctx context.Context) error {
	// An imaginary function that executes a simple database query.
	if err := HealthCheck(ctx); err != nil {
		return err
	}

	return nil
}

func HealthCheck(ctx context.Context) error {
	time.Sleep(500 * time.Millisecond)
	return ctx.Err()
}
