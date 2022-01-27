package main

import (
	"fmt"
	"next/web"
	"testing"
)

func handleHello(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "hello"}
	return ctx.Success(data)
}

func handleRoot(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "root"}
	return ctx.Success(data)
}

func handleAB(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "ab"}
	return ctx.Success(data)
}

func handleABC(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "abc"}
	return ctx.Success(data)
}

func handleTest1(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "test1"}
	panic("i am panic,sad")
	return ctx.Success(data)
}

func handleTest1Test(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "test1test"}
	return ctx.Success(data)
}

func handleTest2(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "test2"}
	return ctx.Success(data)
}

func handleTest2Test(ctx *web.Context) error {
	type Data struct {
		Name string
	}
	data := &Data{Name: "test2test"}
	return ctx.Success(data)
}

func handleLog(ctx *web.Context) error {
	fmt.Println("log")
	return nil
}

func handleAuth(ctx *web.Context) error {
	fmt.Println("auth")
	return nil
}

func TestWeb(t *testing.T) {
	engine := web.NewEngine()
	engine.Get("/hello", handleHello)
	engine.Get("/", handleRoot)
	engine.Get("/a/b", handleAB)
	engine.Get("/a/b/c", handleABC)

	v1 := engine.Group("/v1")
	v1.Use(handleLog)
	v1.Use(handleAuth)
	v1.Get("/test1", handleTest1)
	v1.Get("/test1/test", handleTest1Test)

	v2 := engine.Group("/v2")
	v2.Get("/", handleTest2)
	v2.Get("/test", handleTest2Test)
	engine.Run()
}
