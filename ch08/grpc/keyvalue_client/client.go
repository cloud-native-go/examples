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
	"context"
	"log"
	"os"
	"strings"
	"time"

	pb "github.com/cloud-native-go/ch08/grpc/keyvalue"
	"google.golang.org/grpc"
)

func main() {
	// Use context to establish a 1-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Set up a connection to the gRPC server
	opts := []grpc.DialOption{grpc.WithInsecure(), grpc.WithBlock()}
	conn, err := grpc.DialContext(ctx, "localhost:50051", opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Get a new instance of our client
	client := pb.NewKeyValueClient(conn)

	var action, key, value string

	// Expect something like "set foo bar"
	if len(os.Args) > 2 {
		action, key = os.Args[1], os.Args[2]
		value = strings.Join(os.Args[3:], " ")
	}

	// Call client.Get() or client.Put() as appropriate.
	switch action {
	case "get":
		r, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			log.Fatalf("could not get value for key %s: %v\n", key, err)
		}
		log.Printf("Get %s returns: %s", key, r.Value)

	case "put":
		_, err := client.Put(ctx, &pb.PutRequest{Key: key, Value: value})
		if err != nil {
			log.Fatalf("could not get put key %s: %v\n", key, err)
		}
		log.Printf("Put %s", key)

	default:
		log.Fatalf("Syntax: go run [get|put] KEY VALUE...")
	}
}
