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
	"log"

	"github.com/cloud-native-go/examples/ch08/hexarch/core"
	"github.com/cloud-native-go/examples/ch08/hexarch/frontend"
	"github.com/cloud-native-go/examples/ch08/hexarch/transact"
)

func main() {
	// Create our TransactionLogger. This is an adapter that will plug
	// into the core application's TransactionLogger plug.
	tl, _ := transact.NewTransactionLogger("file")

	// Create Core and tell it which TransactionLogger to use.
	// This is an example of a "driven agent"
	store := core.NewKeyValueStore().WithTransactionLogger(tl)
	store.Restore()

	// Create the frontend.
	// This is an example of a "driving agent".
	fe, _ := frontend.NewFrontEnd("rest")

	log.Fatal(fe.Start(store))
}
