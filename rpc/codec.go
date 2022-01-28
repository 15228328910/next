package rpc

type Codec interface {
	ReadHead(interface{}) error
	ReadBody(interface{}) error
	WriteHead(interface{}) error
	WriteBody(interface{}) error
}

type codec struct {
}

func (c *codec) ReadHead(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *codec) ReadBody(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *codec) WriteHead(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (c *codec) WriteBody(i interface{}) error {
	//TODO implement me
	panic("implement me")
}

func NewCodec() Codec {
	return &codec{}
}
