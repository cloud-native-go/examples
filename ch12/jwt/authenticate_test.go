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
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthenticate(t *testing.T) {
	// Build mock request and a response recorder
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	req.Form = url.Values{
		"username": []string{"foo"},
		"password": []string{"bar"},
	}
	w := httptest.NewRecorder()

	// Call authenticateHandler with mocked values
	authenticateHandler(w, req)

	// Get the recorded response
	resp := w.Result()

	// Just test that we get a 200. Token verification testing
	// is done in TestVerify.
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
