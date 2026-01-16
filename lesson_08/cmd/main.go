package main

import (
	"fmt"

	lru "github.com/pan68ruslan/GoLand/lesson_08/lru"
)

func main() {
	cache := lru.NewLruCache(3)

	cache.Put("A", "Apple")
	cache.Put("B", "Banana")
	cache.Put("C", "Cherry")

	fmt.Println(cache.Get("A"))

	cache.Put("D", "Date")

	fmt.Println(cache.Get("B"))
	fmt.Println(cache.Get("C"))
	fmt.Println(cache.Get("D"))
}
