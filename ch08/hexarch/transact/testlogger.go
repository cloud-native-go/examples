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

package transact

import (
	"github.com/cloud-native-go/examples/ch08/hexarch/core"
)

type TestTransactionLogger struct {
	events   chan<- core.Event // Write-only channel for sending events
	errors   <-chan error      // Read-only channel for receiving errors
	records  []core.Event      // An in-memory record of events
	sequence uint64
}

func (l *TestTransactionLogger) WritePut(key, value string) {
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: value}
}

func (l *TestTransactionLogger) WriteDelete(key string) {
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
}

func (l *TestTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *TestTransactionLogger) LastSequence() uint64 {
	return l.sequence
}

func (l *TestTransactionLogger) Run() {
	l.records = make([]core.Event, 16)

	events := make(chan core.Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() {
		for e := range events { // Retrieve the next Event
			l.records = append(l.records, e)
		}
	}()
}

func (l *TestTransactionLogger) Wait() {}

func (l *TestTransactionLogger) Close() error {
	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return nil
}

func (l *TestTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event) // An unbuffered events channel
	outError := make(chan error, 1)   // A buffered errors channel

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		for _, e := range l.records { // Iterate over the in-memory records
			outEvent <- e // Send e to the channel
		}
	}()

	return outEvent, outError
}

func NewTestTransactionLogger() (core.TransactionLogger, error) {
	return &TestTransactionLogger{}, nil
}
