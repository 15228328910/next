package main

import (
	"fmt"
	"next/rpc"
)

type Test struct {
}

type Resp struct {
	Name string
}

func (t *Test) HelloWorld() (resp *Resp, err error) {
	resp = new(Resp)
	return
}

func main() {
	test := new(Test)
	server := rpc.NewServer()
	server.Register(test)
	server.Run()

	client := rpc.NewClient()
	resp, err := client.Call("Test.HelloWorld", "liucong")
	fmt.Println(resp, err)
}
