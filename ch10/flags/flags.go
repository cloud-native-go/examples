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
	"flag"
	"fmt"
)

func main() {
	// Declare a string flag with a default value "foo"
	// and a short description. It returns a string pointer.
	strp := flag.String("string", "foo", "a string")

	// Declare number and boolean flags, similar to the string flag.
	intp := flag.Int("number", 42, "an integer")
	boolp := flag.Bool("boolean", false, "a boolean")

	// Call flag.Parse() to execute command-line parsing.
	flag.Parse()

	// Print the parsed options and trailing positional arguments.
	fmt.Println("string:", *strp)
	fmt.Println("integer:", *intp)
	fmt.Println("boolean:", *boolp)
	fmt.Println("args:", flag.Args())
}
