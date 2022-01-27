package rpc

type Client interface {
	// Call , invoke service method ,like srv.method
	Call(method string, args interface{}) (reply interface{}, err error)
}

type client struct {
}

func (c *client) Call(method string, args interface{}) (reply interface{}, err error) {
	return
}

func NewClient() Client {
	return &client{}
}
