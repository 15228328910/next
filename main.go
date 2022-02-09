package main

import (
	"fmt"
	"next/web"
)

func main() {
	tire := web.NewTire()
	tire.AddHandler("/a/b/c", nil)
	tire.AddHandler("/a/b", nil)
	tire.Display()
	level := tire.Level()
	fmt.Println("level is :", level)
}
