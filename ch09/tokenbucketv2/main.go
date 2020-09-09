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
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var throttled = Throttle(getHostname, 1, 1, time.Second)

func getHostname(ctx context.Context) (string, error) {
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return os.Hostname()
}

func throttledHandler(w http.ResponseWriter, r *http.Request) {
	ok, hostname, err := throttled(r.Context(), r.RemoteAddr)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !ok {
		http.Error(w, "Too many requests", http.StatusTooManyRequests)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(hostname))
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hostname", throttledHandler)
	log.Fatal(http.ListenAndServe(":8080", r))
}
