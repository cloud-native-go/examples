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
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandlers(t *testing.T) {
	// Build mock request and a response recorder
	req := httptest.NewRequest("GET", "/authenticate", nil)
	req.Form = url.Values{
		"username": []string{"foo"},
		"password": []string{"bar"},
	}
	w := httptest.NewRecorder()

	// Call authenticateHandler with mocked values
	authenticateHandler(w, req)

	// Get the recorded response
	resp := w.Result()
	token, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	// Now we have the token. Let's verify it.

	req2 := httptest.NewRequest("GET", "/verify", nil)
	req2.Header = http.Header{"Authorization": {"Bearer " + string(token)}}
	w2 := httptest.NewRecorder()

	verifyHandler(w2, req2)

	resp2 := w2.Result()
	body, err := io.ReadAll(resp2.Body)
	assert.NoError(t, err)

	assert.Contains(t, string(body), "Welcome back")
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestVerifyNoHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/verify", bytes.NewReader([]byte{}))
	w := httptest.NewRecorder()

	verifyHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}

func TestVerifyInvalidHeader(t *testing.T) {
	req := httptest.NewRequest("GET", "/verify", bytes.NewReader([]byte{}))
	req.Header = http.Header{"Authorization": {"Bearer ABCDEFG"}}
	w := httptest.NewRecorder()

	verifyHandler(w, req)
	resp := w.Result()

	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
}
