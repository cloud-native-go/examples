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

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "cng",
	Long: "A super simple command.",
}

func init() {
	// rootCmd.AddCommand(helloWorldCmd)

	flagsCmd.Flags().StringVarP(&strp, "string", "s", "foo", "a string")
	flagsCmd.Flags().IntVarP(&intp, "number", "n", 42, "an integer")
	flagsCmd.Flags().BoolVarP(&boolp, "boolean", "b", false, "a boolean")
	rootCmd.AddCommand(flagsCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
