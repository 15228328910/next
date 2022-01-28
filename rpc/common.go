package rpc

import "reflect"

//Option  用户协商
type Option struct {
	CodeType int64
}

type Head struct {
	ServiceMethod string
	Seq           string
	Error         string
}

func callFunction(srv interface{}, method string, params interface{}) (ret interface{}, err error) {
	ty := reflect.TypeOf(srv)
	m, _ := ty.MethodByName(method)
	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(srv)
	in[1] = reflect.ValueOf(params)
	values := m.Func.Call(in)
	ret = values[0].Interface()
	if values[1].Interface() == nil {
		err = nil
	} else {
		err = values[1].Interface().(error)
	}
	return
}
