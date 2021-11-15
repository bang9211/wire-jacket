package mockup

import viperjacket "github.com/bang9211/viper-jacket"

type ExplorerServer interface {
	// Serve serves server.
	Serve() error
	// GetAllBlockData return all blocks of blockchain.
	GetAllBlockData() []string
	// Close closes blockchain.
	Close() error
}

type MockupExplorerServer struct {
	config     viperjacket.Config
	blockchain Blockchain
}

func NewMockupExplorerServer(config viperjacket.Config, blockchain Blockchain) ExplorerServer {
	return &MockupExplorerServer{config: config, blockchain: blockchain}
}

func (mes *MockupExplorerServer) Serve() error {
	return nil
}

func (mes *MockupExplorerServer) GetAllBlockData() []string {
	allBlockData := []string{}
	for _, block := range mes.blockchain.GetBlocks() {
		allBlockData = append(allBlockData, block.GetData())
	}
	return allBlockData
}

func (mes *MockupExplorerServer) Close() error {
	return nil
}
