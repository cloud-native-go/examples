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
	"database/sql"
	"fmt"
	"net/url"
	"sync"

	"github.com/cloud-native-go/examples/ch08/hexarch/core"
	_ "github.com/lib/pq" // Load the Postgres drivers
)

type PostgresDbParams struct {
	dbName   string
	host     string
	user     string
	password string
}

type PostgresTransactionLogger struct {
	events chan<- core.Event // Write-only channel for sending events
	errors <-chan error      // Read-only channel for receiving errors
	db     *sql.DB           // Our database access interface
	wg     *sync.WaitGroup   // Used to ensure writes are completed
}

func (l *PostgresTransactionLogger) WritePut(key, value string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventPut, Key: key, Value: url.QueryEscape(value)}
	l.wg.Done()
}

func (l *PostgresTransactionLogger) WriteDelete(key string) {
	l.wg.Add(1)
	l.events <- core.Event{EventType: core.EventDelete, Key: key}
	l.wg.Done()
}

func (l *PostgresTransactionLogger) Err() <-chan error {
	return l.errors
}

func (l *PostgresTransactionLogger) LastSequence() uint64 {
	return 0
}

func (l *PostgresTransactionLogger) Run() {
	events := make(chan core.Event, 16) // Make an events channel
	l.events = events

	errors := make(chan error, 1) // Make an errors channel
	l.errors = errors

	go func() { // The INSERT query
		query := `INSERT INTO transactions
			(event_type, key, value)
			VALUES ($1, $2, $3)`

		for e := range events { // Retrieve the next Event
			_, err := l.db.Exec( // Execute the INSERT query
				query,
				e.EventType, e.Key, e.Value)

			if err != nil {
				errors <- err
			}
		}
	}()
}

func (l *PostgresTransactionLogger) Wait() {
	l.wg.Wait()
}

func (l *PostgresTransactionLogger) Close() error {
	l.wg.Wait()

	if l.events != nil {
		close(l.events) // Terminates Run loop and goroutine
	}

	return l.db.Close()
}

func (l *PostgresTransactionLogger) ReadEvents() (<-chan core.Event, <-chan error) {
	outEvent := make(chan core.Event) // An unbuffered events channel
	outError := make(chan error, 1)   // A buffered errors channel

	query := "SELECT sequence, event_type, key, value FROM transactions"

	go func() {
		defer close(outEvent) // Close the channels when the
		defer close(outError) // goroutine ends

		rows, err := l.db.Query(query) // Run query; get result set
		if err != nil {
			outError <- fmt.Errorf("sql query error: %w", err)
			return
		}

		defer rows.Close() // This is important!

		var e core.Event // Create an empty Event

		for rows.Next() { // Iterate over the rows

			err = rows.Scan( // Read the values from the
				&e.Sequence, &e.EventType, // row into the Event.
				&e.Key, &e.Value)

			if err != nil {
				outError <- err
				return
			}

			outEvent <- e // Send e to the channel
		}

		err = rows.Err()
		if err != nil {
			outError <- fmt.Errorf("transaction log read failure: %w", err)
		}
	}()

	return outEvent, outError
}

func (l *PostgresTransactionLogger) verifyTableExists() (bool, error) {
	const table = "transactions"

	var result string

	rows, err := l.db.Query(fmt.Sprintf("SELECT to_regclass('public.%s');", table))
	defer rows.Close()
	if err != nil {
		return false, err
	}

	for rows.Next() && result != table {
		rows.Scan(&result)
	}

	return result == table, rows.Err()
}

func (l *PostgresTransactionLogger) createTable() error {
	var err error

	createQuery := `CREATE TABLE transactions (
		sequence      BIGSERIAL PRIMARY KEY,
		event_type    SMALLINT,
		key 		  TEXT,
		value         TEXT
	  );`

	_, err = l.db.Exec(createQuery)
	if err != nil {
		return err
	}

	return nil
}

func NewPostgresTransactionLogger(param PostgresDbParams) (core.TransactionLogger, error) {
	connStr := fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable",
		param.host, param.dbName, param.user, param.password)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to create db value: %w", err)
	}

	err = db.Ping() // Test the databases connection
	if err != nil {
		return nil, fmt.Errorf("failed to opendb connection: %w", err)
	}

	tl := &PostgresTransactionLogger{db: db, wg: &sync.WaitGroup{}}

	exists, err := tl.verifyTableExists()
	if err != nil {
		return nil, fmt.Errorf("failed to verify table exists: %w", err)
	}
	if !exists {
		if err = tl.createTable(); err != nil {
			return nil, fmt.Errorf("failed to create table: %w", err)
		}
	}

	return tl, nil
}
