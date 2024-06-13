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
	"os"
	"testing"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func TestCreateLogger(t *testing.T) {
	const filename = "/tmp/create-logger.txt"
	t.Cleanup(func() {
		err := os.Remove(filename)
		if err != nil {
			return
		}
	})

	tl, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Errorf("Got error: %v", err)
	}
	if tl == nil {
		t.Error("Logger is nil?")
	}

	if !fileExists(filename) {
		t.Errorf("File %s doesn't exist", filename)
	}
}

func TestLastSequence(t *testing.T) {
	const filename = "/tmp/last-sequence.txt"
	t.Cleanup(func() {
		err := os.Remove(filename)
		if err != nil {
			return
		}
	})

	tl, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Error(err)
	}
	tl.Run()
	t.Cleanup(func() {
		closeErr := tl.Close()
		if closeErr != nil {
			return
		}
	})

	evaluateLastSequence(t, tl, 0)

	tl.WritePut("my-key", "my-value")
	tl.Wait()

	evaluateLastSequence(t, tl, 1)

	tl.WritePut("my-key", "my-value")
	tl.WritePut("my-key", "my-value2")
	tl.Wait()

	evaluateLastSequence(t, tl, 3)
}

func TestWriteAppend(t *testing.T) {
	const filename = "/tmp/write-append.txt"
	t.Cleanup(func() {
		err := os.Remove(filename)
		if err != nil {
			return
		}
	})

	tl, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Error(err)
	}
	tl.Run()
	t.Cleanup(func() {
		closeErr := tl.Close()
		if closeErr != nil {
			return
		}
	})

	chev, cherr := tl.ReadEvents()
	for e := range chev {
		t.Log(e)
	}
	err = <-cherr
	if err != nil {
		t.Error(err)
	}

	tl.WritePut("my-key", "my-value")
	tl.WritePut("my-key", "my-value2")
	tl.Wait()

	tl2, err := NewFileTransactionLogger(filename)
	if err != nil {
		t.Error(err)
	}
	tl2.Run()
	defer tl2.Close()

	chev, cherr = tl2.ReadEvents()
	for e := range chev {
		t.Log(e)
	}
	err = <-cherr
	if err != nil {
		t.Error(err)
	}

	tl2.WritePut("my-key", "my-value3")
	tl2.WritePut("my-key2", "my-value4")
	tl2.Wait()

	if tl2.LastSequence() != 4 {
		t.Errorf("Last sequence mismatch (expected 4; got %d)", tl2.LastSequence())
	}
}

func TestWritePut(t *testing.T) {
	const filename = "/tmp/write-put.txt"
	t.Cleanup(func() {
		err := os.Remove(filename)
		if err != nil {
			return
		}
	})

	tl, _ := NewFileTransactionLogger(filename)
	tl.Run()
	t.Cleanup(func() {
		closeErr := tl.Close()
		if closeErr != nil {
			return
		}
	})

	tl.WritePut("my-key", "my-value")
	tl.WritePut("my-key", "my-value2")
	tl.WritePut("my-key", "my-value3")
	tl.WritePut("my-key", "my-value4")
	tl.Wait()

	tl2, _ := NewFileTransactionLogger(filename)
	evin, errin := tl2.ReadEvents()
	defer tl2.Close()

	for e := range evin {
		t.Log(e)
	}

	err := <-errin
	if err != nil {
		t.Error(err)
	}

	if tl.LastSequence() != tl2.LastSequence() {
		t.Errorf("Last sequence mismatch (%d vs %d)", tl.LastSequence(), tl2.LastSequence())
	}
}

func evaluateLastSequence(t *testing.T, el TransactionLogger, expected uint64) {
	if ls := el.LastSequence(); ls != expected {
		t.Errorf("Last sequence mismatch (expected %d; got %d)", expected, ls)
	} else {
		t.Logf("Last sequence agrees with expectations (expected %d; got %d)", expected, ls)
	}
}
