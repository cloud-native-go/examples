/*
 * Copyright 2024 Matthew A. Titmus
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

package ch04

import (
	"hash/fnv"
	"reflect"
	"sync"
)

type Shard[K comparable, V any] struct {
	sync.RWMutex         // Compose from sync.RWMutex
	items        map[K]V // m contains the shard's data
}

type ShardedMap[K comparable, V any] []*Shard[K, V]

// NewShardedMap creates and initializes a new ShardedMap with the specified
// number of shards.
func NewShardedMap[K comparable, V any](nshards int) ShardedMap[K, V] {
	shards := make([]*Shard[K, V], nshards) // Initialize a *Shards slice

	for i := 0; i < nshards; i++ {
		shard := make(map[K]V)
		shards[i] = &Shard[K, V]{items: shard} // A ShardedMap IS a slice!
	}

	return shards
}

// getShardIndex accepts a key and returns a value in 0..N-1, where N is
// the number of shards.
func (m ShardedMap[K, V]) getShardIndex(key K) int {
	str := reflect.ValueOf(key).String() // Get string representation of key
	hash := fnv.New32a()                 // Get a hash implementation from "hash/fnv"
	hash.Write([]byte(str))              // Write bytes to the hash
	sum := int(hash.Sum32())             // Get the resulting checksum
	return sum % len(m)                  // Mod by len(m) to get index
}

// getShard accepts a key and returns a pointer to its corresponding Shard.
func (m ShardedMap[K, V]) getShard(key K) *Shard[K, V] {
	index := m.getShardIndex(key)
	return m[index]
}

// Delete removes a value from the map. If key doesn't exist in the map,
// this method is a no-op.
func (m ShardedMap[K, V]) Delete(key K) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.items, key)
}

// Get retrieves and returns a value from the map. If the value doesn't exist,
// nil is returned.
func (m ShardedMap[K, V]) Get(key K) V {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.items[key]
}

func (m ShardedMap[K, V]) Set(key K, value V) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.items[key] = value
}

// Keys returns a list of all keys in the sharded map.
func (m ShardedMap[K, V]) Keys() []K {
	var keys []K         // Declare an empty keys slice
	var mutex sync.Mutex // Mutex for write safety to keys

	var wg sync.WaitGroup // Create a wait group and add a
	wg.Add(len(m))        // wait value for each slice

	for _, shard := range m { // Run a goroutine for each slice in m
		go func(s *Shard[K, V]) {
			s.RLock() // Establish a read lock on s

			defer wg.Done()   // Release of the read lock
			defer s.RUnlock() // Tell the WaitGroup it's done

			for key, _ := range s.items { // Get the slice's keys
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}
		}(shard)
	}

	wg.Wait() // Block until all goroutines are done

	return keys // Return combined keys slice
}
