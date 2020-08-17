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

package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// Sayer says what an animal says.
type Sayer interface {
	Says() string
}

// Here is an implementation that talks over RPC
type SayerRPC struct{ client *rpc.Client }

func (g *SayerRPC) Says() string {
	var resp string
	err := g.client.Call("Plugin.Says", new(interface{}), &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that SayerRPC talks to, conforming to
// the requirements of net/rpc
type SayerRPCServer struct {
	// This is the real implementation
	Impl Sayer
}

func (s *SayerRPCServer) Says(args interface{}, resp *string) error {
	*resp = s.Impl.Says()
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a SayerRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return SayerRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type SayerPlugin struct {
	Impl Sayer
}

func (p *SayerPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &SayerRPCServer{Impl: p.Impl}, nil
}

func (SayerPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &SayerRPC{client: c}, nil
}

// HandshakeConfig is a common handshake that is shared by plugin and host.
var HandshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}
