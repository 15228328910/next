package web

import (
	"fmt"
	"log"
)

//Recover 默认崩溃恢复中间件
func Recover(ctx *Context) error {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recover from panic")
		}
	}()
	return nil
}

//Log 默认日志中间件
func Log(ctx *Context) error {
	log.Println("打印日志.......")
	return nil
}
