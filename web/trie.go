package web

import (
	"fmt"
	"strings"
	"unsafe"
)

type Tire struct {
	children []*Tire
	// 当前节点
	hold string
	// 全路径
	path string
	// 处理函数
	handler Service
}

var tire Tire

func NewTire() *Tire {
	return &tire
}

func (t *Tire) AddNode(hold string, handler Service) *Tire {
	if hold == "" {
		return t
	}
	for _, child := range t.children {
		// match
		if child.hold == hold {
			return child
		}
	}

	path := "/" + hold
	if !t.isRootNode() {
		path = t.path + "/" + hold
	}
	// not match
	child := &Tire{
		children: nil,
		hold:     hold,
		path:     path,
		handler:  handler,
	}
	t.children = append(t.children, child)
	return child
}

func (t *Tire) isRootNode() bool {
	return uintptr(unsafe.Pointer(t)) == uintptr(unsafe.Pointer(&tire))
}

func (t *Tire) AddHandler(path string, handler Service) *Tire {
	path = strings.Trim(path, "/")
	if path == "" {
		// 如果是首节点
		if t.isRootNode() {
			t.handler = handler
			t.path = "/"
		}
		return t
	}
	pathArray := strings.SplitN(path, "/", 2)
	child := t.AddNode(pathArray[0], handler)
	if len(pathArray) == 2 {
		return child.AddHandler(pathArray[1], handler)
	}
	return t
}

//GetHandler 匹配路由
func (t *Tire) GetHandler(path string) Service {
	if t.path == path {
		return t.handler
	}
	for _, child := range t.children {
		if handler := child.GetHandler(path); handler != nil {
			return handler
		}
	}
	return nil
}

func (t *Tire) Display() {
	if t.children == nil {
		fmt.Println(t.path)
		return
	}
	for _, child := range t.children {
		child.Display()
	}
}
