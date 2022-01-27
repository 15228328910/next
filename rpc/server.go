package rpc

type Server interface {
	// Register all exported method of service
	Register(srv Service) error
	Run() error
}

type server struct {
}

func (s *server) Register(srv Service) error {
	return nil
}

func (s *server) Run() error {
	return nil
}

func NewServer() Server {
	return &server{}
}
