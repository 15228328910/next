package rpc

import (
	"encoding/gob"
	"io"
)

type gobCodec struct {
	encoder *gob.Encoder
	decoder *gob.Decoder
	closer  io.ReadWriteCloser
}

func (c *gobCodec) ReadHead(i interface{}) error {
	return c.decoder.Decode(i)
}

func (c *gobCodec) ReadBody(i interface{}) error {
	return c.decoder.Decode(i)
}

func (c *gobCodec) WriteHead(i interface{}) error {
	return c.encoder.Encode(i)
}

func (c *gobCodec) WriteBody(i interface{}) error {
	return c.encoder.Encode(i)
}

func newGobCodec(closer io.ReadWriteCloser) Codec {
	return &gobCodec{
		closer:  closer,
		encoder: gob.NewEncoder(closer),
		decoder: gob.NewDecoder(closer),
	}
}
