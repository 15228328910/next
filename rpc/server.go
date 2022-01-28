package rpc

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"
)

type Server interface {
	// Register all exported method of service
	Register(srv interface{})
	Run() error
}

type server struct {
	srvMap map[string]interface{}
	addr   string
}

func (s *server) Register(srv interface{}) {
	ty := reflect.TypeOf(srv)
	name := ty.String()
	if _, ok := s.srvMap[name]; ok {
		panic("服务" + name + "重复注册")
	}
	s.srvMap[name] = srv
}

func (s *server) Run() error {
	return s.run()
}

func (s *server) run() error {
	listen, err := net.Listen("tcp", s.addr)
	if err != nil {
		panic(err)
	}
	fmt.Println("running......")
	for {
		fmt.Println("waiting for accept......")
		conn, errConn := listen.Accept()
		if errConn != nil {
			fmt.Println("connect error", errConn)
		}
		fmt.Println("waiting for handle......")
		go s.handleConn(conn)
	}
}

func (s *server) handleConn(conn io.ReadWriteCloser) {
	var option Option
	json.NewDecoder(conn).Decode(&option)
	fmt.Println("codecType is:", option.CodeType)

	decoder := NewCodecFactory(option.CodeType).GetCodec(conn)
	// read header
	header := new(Head)
	err := decoder.ReadHead(header)
	if err != nil {
		header.Error = err.Error()
	}

	// read body
	var body interface{}
	err = decoder.ReadBody(&body)
	if err != nil {
		header.Error = err.Error()
	}

	// 获取注册方法
	srvMethodArr := strings.Split(header.ServiceMethod, "/")
	srv := s.srvMap[srvMethodArr[0]]
	method := srvMethodArr[1]
	resp, err := callFunction(srv, method, body)
	if err != nil {
		header.Error = err.Error()
	}

	// 写入头部
	encoder := json.NewEncoder(conn)
	encoder.Encode(header)
	encoder.Encode(resp)
}

func NewServer(addr string) Server {
	return &server{
		srvMap: make(map[string]interface{}),
		addr:   addr,
	}
}
