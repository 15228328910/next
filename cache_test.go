package main

import (
	"fmt"
	"next/cache"
	"testing"
)

func TestCache(t *testing.T) {
	c := cache.NewCache(2)
	c.Set("c", "c")
	c.Set("a", "a")
	c.Set("b", "b")
	fmt.Println("b", c.Get("b"))
	fmt.Println("a", c.Get("a"))
	fmt.Println("c", c.Get("c"))
	c.Set("a", 1)
	c.Set("b", 2)
	fmt.Println("b", c.Get("b"))
	fmt.Println("a", c.Get("a"))

}
