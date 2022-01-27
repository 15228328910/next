package rpc

type Client interface {
	// Call , invoke service method ,like srv.method
	Call(method string, args interface{}) (reply interface{}, err error)
}
