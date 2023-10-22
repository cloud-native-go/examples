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
	"github.com/cloud-native-go/ch08/hashicorp-plugin/commons"
	"github.com/hashicorp/go-plugin"
)

// Here is a real implementation of Sayer
type Duck struct{}

func (g *Duck) Says() string {
	return "Quack!"
}

func main() {
	sayer := &Duck{}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"sayer": &commons.SayerPlugin{Impl: sayer},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: commons.HandshakeConfig,
		Plugins:         pluginMap,
	})
}
