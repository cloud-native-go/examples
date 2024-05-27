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
	"fmt"
	"log"
	"net/http"
	"strings"
)

func verifyHandler(w http.ResponseWriter, r *http.Request) {
	header := r.Header.Get("Authorization")
	token := strings.TrimPrefix(header, "Bearer ")

	if header == "" || token == header {
		log.Println("Missing authorization header")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	claims, err := verifyToken(token)
	if err != nil {
		log.Println("Invalid authorization token:", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	fmt.Fprintf(w, "Welcome back, %s!\n", claims.Name)
}
