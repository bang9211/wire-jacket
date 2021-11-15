package mockup

import viperjacket "github.com/bang9211/viper-jacket"

type RESTAPIServer interface {
	// Serve serves server.
	Serve() error
	// GetPaths gets all the paths of REST API.
	GetPaths() []string
	// Close closes blockchain.
	Close() error
}

type MockupRESTAPIServer struct {
	config     viperjacket.Config
	blockchain Blockchain
}

func NewMockupRESTAPIServer(config viperjacket.Config, blockchain Blockchain) RESTAPIServer {
	return &MockupRESTAPIServer{config: config, blockchain: blockchain}
}

func (mrs *MockupRESTAPIServer) Serve() error {
	return nil
}

func (mrs *MockupRESTAPIServer) GetPaths() []string {
	return []string{"/"}
}

func (mrs *MockupRESTAPIServer) Close() error {
	return nil
}
