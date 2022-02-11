package main

import (
	"fmt"
	"next/cache"
	"testing"
)

func TestRing(t *testing.T) {
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
