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

package main

import (
	"fmt"
	"os"
	"plugin"
)

// Sayer says what an animal says.
type Sayer interface {
	Says() string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: run main/main.go animal")
		os.Exit(1)
	}

	// Get the animal name, and build the path where we expect to
	// find the corresponding shared object (.so) file.
	name := os.Args[1]
	module := fmt.Sprintf("./%s/%s.so", name, name)

	// Does the file exist?
	_, err := os.Stat(module)
	if os.IsNotExist(err) {
		fmt.Println("can't find an animal named", name)
		os.Exit(1)
	}

	// Open our plugin. and returns a *plugin.Plugin.
	p, err := plugin.Open(module)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Lookup searches for a symbol, which can be any exported variable
	// or function, named "Animal" in plugin p.
	symbol, err := p.Lookup("Animal")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Asserts that the symbol interface holds an Sayer.
	animal, ok := symbol.(Sayer)
	if !ok {
		fmt.Println("that's not an Sayer")
		os.Exit(1)
	}

	// Now we can use our loaded plugin!
	fmt.Printf("A %s says: %q\n", name, animal.Says())
}
