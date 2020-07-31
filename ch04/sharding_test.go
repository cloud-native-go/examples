package ch04

import (
	"testing"
)

// TestShardingGetShardIndex tests getShardIndex. It's not a very good test: it just
// gets the indices of 5 keys and fails the test if there are any hash
// collisions, so it could fail if you change the hash algorithm.
func TestShardingGetShardIndex(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)
	counts := make([]int, BUCKETS)

	keys := []string{"A", "B", "C", "D", "E"}

	for _, key := range keys {
		idx := sMap.getShardIndex(key)
		counts[idx]++
		t.Log(key, idx)
		if counts[idx] > 1 {
			t.Error("Hash collision")
		}
	}
}

// TestShardingSetAndGet tests Set and Get... by setting some values and getting them.
func TestShardingSetAndGet(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	for k, v := range truthMap {
		got := sMap.Get(k)

		if got != v {
			t.Errorf("Key mismatch on %s: expected %d, got %d", k, v, got)
		}
	}
}

// TestShardingKeys tests the Keys method by adding 5 values to the map and checking
// that each one exists in the keys list exactly once.
func TestShardingKeys(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	keys := sMap.Keys()

	if len(truthMap) != len(keys) {
		t.Error("Map/keys mismatch")
	}

	for _, key := range sMap.Keys() {
		if _, ok := truthMap[key]; !ok {
			t.Error("Key", key, "not in truthMap")
		}

		delete(truthMap, key)
	}

	if len(truthMap) != 0 {
		t.Error("Key mismatch")
	}
}

// TestShardingDelete tests the Delete method by adding and then removing five values.
func TestShardingDelete(t *testing.T) {
	const BUCKETS = 17

	sMap := NewShardedMap(BUCKETS)

	truthMap := map[string]int{
		"alpha":   1,
		"beta":    2,
		"gamma":   3,
		"delta":   4,
		"epsilon": 5,
	}

	for k, v := range truthMap {
		sMap.Set(k, v)
	}

	keys := sMap.Keys()
	for _, key := range keys {
		sMap.Delete(key)
	}

	if len(sMap.Keys()) != 0 {
		t.Error("Deletion failure")
	}
}
