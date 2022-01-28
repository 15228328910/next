package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"reflect"
)

type Client interface {
	// Call , invoke service method ,like srv.method
	Call(object interface{}, method string, args interface{}) (reply interface{}, err error)
}

type client struct {
	addr         string
	codecFactory CodecFactory
	codecType    int64
}

func (c *client) Call(object interface{}, method string, args interface{}) (reply interface{}, err error) {
	fmt.Print("准备连接")
	conn, err := net.Dial("tcp", c.addr)
	if err != nil {
		return
	}

	fmt.Println("准备发起请求数据,codeType:", c.codecType)
	json.NewEncoder(conn).Encode(&Option{CodeType: 1})

	target := reflect.TypeOf(object).String()
	// 写入头部
	encoder := c.codecFactory.GetCodec(conn)
	header := &Head{
		ServiceMethod: target + "/" + method,
		Seq:           "1",
	}
	encoder.WriteHead(header)

	// 写入body
	encoder.WriteBody(args)

	// 读取头部
	decoder := json.NewDecoder(conn)
	decoder.Decode(header)
	err = errors.New(header.Error)

	// 读取body
	decoder.Decode(&reply)
	return
}

func NewClient(addr string, codecType int64) Client {
	return &client{addr: addr, codecFactory: NewCodecFactory(codecType), codecType: codecType}
}
