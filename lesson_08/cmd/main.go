package main

import (
	"fmt"
	lru "lesson_08/lru"
)

func main() {
	cache := lru.NewLruCache(3)

	cache.Put("A", "Apple")
	cache.Put("B", "Banana")
	cache.Put("C", "Cherry")

	fmt.Println(cache.Get("A")) // Apple, true

	cache.Put("D", "Date") // витіснить найстаріший

	fmt.Println(cache.Get("B")) // "", false (видалений)
	fmt.Println(cache.Get("C")) // Cherry, true (ще є)
	fmt.Println(cache.Get("D")) // Date, true
}
