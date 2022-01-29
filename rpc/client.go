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
	Close() error
}

type client struct {
	addr         string
	codecFactory CodecFactory
	codecType    int64
	conn         net.Conn
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Call(object interface{}, method string, args interface{}) (reply interface{}, err error) {
	fmt.Println("准备发起请求数据,codeType:", c.codecType)
	json.NewEncoder(c.conn).Encode(&Option{CodeType: c.codecType})

	target := reflect.TypeOf(object).String()
	// 写入头部
	encoder := c.codecFactory.GetCodec(c.conn)
	header := &Head{
		ServiceMethod: target + "/" + method,
		Seq:           "1",
	}
	encoder.WriteHead(header)

	// 写入body
	encoder.WriteBody(args)

	// 读取头部
	decoder := json.NewDecoder(c.conn)
	decoder.Decode(header)
	err = errors.New(header.Error)

	// 读取body
	decoder.Decode(&reply)
	return
}

func NewClient(addr string, codecType int64) Client {
	fmt.Print("准备连接")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		fmt.Println("连接失败")
	}
	cli := &client{addr: addr, codecFactory: NewCodecFactory(codecType), codecType: codecType, conn: conn}
	return cli
}
