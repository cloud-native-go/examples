package main

import (
	"fmt"

	lru "github.com/hashicorp/golang-lru/v2"
)

var cache *lru.Cache[int, string]

func init() {
	cache, _ = lru.NewWithEvict(2,
		func(key int, value string) {
			fmt.Printf("Evicted: key=%d value=%s\n", key, value)
		},
	)
}

func main() {
	cache.Add(1, "a") // adds 1
	cache.Add(2, "b") // adds 2; cache is now at capacity

	fmt.Println(cache.Get(1)) // "a true"; 1 now most recently used

	cache.Add(3, "c") // adds 3, evicts key 2

	fmt.Println(cache.Get(2)) // " false" (not found)
}
