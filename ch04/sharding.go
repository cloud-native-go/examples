package ch04

import (
	"crypto/sha1"
	"sync"
)

type Shard struct {
	sync.RWMutex
	m map[string]interface{}
}

type ShardedMap []*Shard

// NewShardedMap creates and initializes a new ShardedMap with the specified
// number of shards.
func NewShardedMap(nshards int) ShardedMap {
	shards := make([]*Shard, nshards)

	for i := 0; i < nshards; i++ {
		shard := make(map[string]interface{})
		shards[i] = &Shard{m: shard}
	}

	return shards
}

// getShardIndex accepts a key and returns a value in 0..N-1, where N is
// the number of shards. As currently written the hash algorithm only works
// correctly for up to 255 shards.
func (m ShardedMap) getShardIndex(key string) int {
	hash := sha1.Sum([]byte(key))

	// Grab an arbitrary byte and mod it by the number of shards
	return int(hash[17]) % len(m)
}

// getShard accepts a key and returns a pointer to its corresponding Shard.
func (m ShardedMap) getShard(key string) *Shard {
	index := m.getShardIndex(key)
	return m[index]
}

// Delete removes a value from the map. If key doesn't exist in the map,
// this method is a no-op.
func (m ShardedMap) Delete(key string) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	delete(shard.m, key)
}

// Get retrieves and returns a value from the map. If the value doesn't exist,
// nil is returned.
func (m ShardedMap) Get(key string) interface{} {
	shard := m.getShard(key)
	shard.RLock()
	defer shard.RUnlock()

	return shard.m[key]
}

func (m ShardedMap) Set(key string, value interface{}) {
	shard := m.getShard(key)
	shard.Lock()
	defer shard.Unlock()

	shard.m[key] = value
}

// Keys returns a list of all keys in the sharded map.
func (m ShardedMap) Keys() []string {
	keys := make([]string, 0)
	mutex := sync.Mutex{}

	wg := sync.WaitGroup{}
	wg.Add(len(m))

	for _, shard := range m {
		go func(s *Shard) {
			s.RLock()

			for key := range s.m {
				mutex.Lock()
				keys = append(keys, key)
				mutex.Unlock()
			}

			s.RUnlock()
			wg.Done()
		}(shard)
	}

	wg.Wait() // Block until all reads are done

	return keys // Return combined keys slice
}
