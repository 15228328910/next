package rpc

import "io"

type Codec interface {
	ReadHead(interface{}) error
	ReadBody(interface{}) error
	WriteHead(interface{}) error
	WriteBody(interface{}) error
}

type CodecFactory interface {
	GetCodec(closer io.ReadWriteCloser) Codec
}

func NewCodecFactory(codecType int64) CodecFactory {
	switch codecType {
	case 1:
		return &jsonCodecFactory{}
	case 2:
		return &gobCodecFactory{}
	}
	return &gobCodecFactory{}
}

type jsonCodecFactory struct {
}

func (j *jsonCodecFactory) GetCodec(closer io.ReadWriteCloser) Codec {
	return newJsonCodec(closer)
}

type gobCodecFactory struct {
}

func (g *gobCodecFactory) GetCodec(closer io.ReadWriteCloser) Codec {
	return newGobCodec(closer)
}
