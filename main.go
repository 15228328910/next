package main

import (
	"fmt"

	"github.com/15228328910/next/cache"
)

func main() {
	r := cache.NewRing()
	r.Add("1")
	r.Add("2")
	r.Add("3")

	mp := map[string]int{
		"1": 0,
		"2": 0,
		"3": 0,
	}
	for i := 0; i < 100; i++ {
		nodeName := fmt.Sprintf("%d", i)
		node := r.Get(nodeName)
		mp[node]++
	}
	for key, value := range mp {
		fmt.Println(key, value)
	}
}
