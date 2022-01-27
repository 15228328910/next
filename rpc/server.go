package rpc

import "reflect"

type Server interface {
	// Register all exported method of service
	Register(srv interface{})
	Run() error
}

type server struct {
	srvMap map[string]interface{}
}

func (s *server) Register(srv interface{}) {
	ty := reflect.TypeOf(srv)
	name := ty.Name()
	if _, ok := s.srvMap[name]; ok {
		panic("服务" + name + "重复注册")
	}
	s.srvMap[name] = srv
}

func (s *server) Run() error {
	return nil
}

func NewServer() Server {
	return &server{
		srvMap: make(map[string]interface{}),
	}
}
