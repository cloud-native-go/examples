/*
 * Copyright 2023 Matthew A. Titmus
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

	"github.com/go-errors/errors"
)

func TestPut(t *testing.T) {
	const key = "create-key"
	const value = "create-value"

	var val interface{}
	var contains bool

	defer delete(store.m, key)

	// Sanity check
	_, contains = store.m[key]
	if contains {
		t.Error("key/value already exists")
	}

	// err should be nil
	err := Put(key, value)
	if err != nil {
		t.Error(err)
	}

	val, contains = store.m[key]
	if !contains {
		t.Error("create failed")
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestGet(t *testing.T) {
	const key = "read-key"
	const value = "read-value"

	var val interface{}
	var err error

	defer delete(store.m, key)

	// Read a non-thing
	val, err = Get(key)
	if err == nil {
		t.Error("expected an error")
	}
	if !errors.Is(err, ErrorNoSuchKey) {
		t.Error("unexpected error:", err)
	}

	store.m[key] = value

	val, err = Get(key)
	if err != nil {
		t.Error("unexpected error:", err)
	}

	if val != value {
		t.Error("val/value mismatch")
	}
}

func TestDelete(t *testing.T) {
	const key = "delete-key"
	const value = "delete-value"

	var contains bool

	defer delete(store.m, key)

	store.m[key] = value

	_, contains = store.m[key]
	if !contains {
		t.Error("key/value doesn't exist")
	}

	Delete(key)

	_, contains = store.m[key]
	if contains {
		t.Error("Delete failed")
	}
}
