package main

import (
	"fmt"
	"next/cache"
)

func main() {
	r := cache.NewRing()
	r.Add("1")
	r.Add("2")
	r.Add("3")

	node := r.Get("2")
	fmt.Println(node)

	r.Remove("2")
	node = r.Get("2")
	fmt.Println(node)

	r.Remove("3")
	node = r.Get("2")
	fmt.Println(node)
}
