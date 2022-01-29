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

func (t *Test) HelloWorld(name string) (resp *Resp, err error) {
	resp = &Resp{
		Name: name,
	}
	return
}

func main() {
	test := new(Test)
	server := rpc.NewServer("127.0.0.1:8096")
	server.Register(test)
	go server.Run()

	client := rpc.NewClient("127.0.0.1:8096", 2)
	defer client.Close()
	resp, err := client.Call(test, "HelloWorld", "liucong")
	fmt.Println(resp, err)
}
