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
	"net/http"
)

func authenticateHandler(w http.ResponseWriter, r *http.Request) {
	// This is required to populate r.Form.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Retrieve and validate the POST-ed credentials
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Authenticate the password, responding to errors appropriately
	valid, err := authenticatePassword(username, password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	} else if !valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Password is valid; build a new token
	tokenString, err := buildToken(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Respond with the new token string
	fmt.Fprint(w, tokenString)
}

// authenticatePassword always returns true and a nil error, just
// for the sake of demonstration
func authenticatePassword(_, _ string) (bool, error) {
	return true, nil
}
