package rpc

type Server interface {
	// Register all exported method of service
	Register(srv Service) error
	// Call , invoke service method ,like srv.method
	Call(method string, args interface{}) (reply interface{}, err error)
}
