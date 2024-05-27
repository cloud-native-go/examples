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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAll(t *testing.T) {
	// 32 bytes for AES-256. Obviously, don't do this.
	key := []byte("example.key.12345678.example.key")

	plaintext := []byte("Hello, Cloud Native Go!")

	encrypted, err := encryptAES(key, plaintext)
	assert.NoError(t, err)

	decrypted, err := decryptAES(key, encrypted)
	assert.NoError(t, err)

	assert.Equal(t, plaintext, decrypted)
}
