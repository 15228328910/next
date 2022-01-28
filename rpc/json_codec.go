package rpc

import (
	"encoding/json"
	"io"
)

type jsonCodec struct {
	encoder *json.Encoder
	decoder *json.Decoder
	closer  io.ReadWriteCloser
}

func (c *jsonCodec) ReadHead(i interface{}) error {
	return c.decoder.Decode(i)
}

func (c *jsonCodec) ReadBody(i interface{}) error {
	return c.decoder.Decode(i)
}

func (c *jsonCodec) WriteHead(i interface{}) error {
	return c.encoder.Encode(i)
}

func (c *jsonCodec) WriteBody(i interface{}) error {
	return c.encoder.Encode(i)
}

func newJsonCodec(closer io.ReadWriteCloser) Codec {
	return &jsonCodec{
		closer:  closer,
		encoder: json.NewEncoder(closer),
		decoder: json.NewDecoder(closer),
	}
}
