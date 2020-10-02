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

	"github.com/spf13/cobra"
)

var strp string
var intp int
var boolp bool

var flagsCmd = &cobra.Command{
	Use:   "flags",
	Short: "Experiment with flags",
	Long:  "A simple flags experimentation command, build with Cobra.",
	Run:   flagsFunc,
}

func flagsFunc(cmd *cobra.Command, args []string) {
	fmt.Println("string:", strp)
	fmt.Println("integer:", intp)
	fmt.Println("boolean:", boolp)
	fmt.Println("args:", args)
}
