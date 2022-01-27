package rpc

import "net/rpc"

type Server interface {
	rpc.Call
}
