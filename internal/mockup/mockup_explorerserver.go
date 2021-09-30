package mockup

import "github.com/bang9211/wire-jacket/internal/config"

type ExplorerServer interface {
	// Serve serves server.
	Serve() error
	// GetAllBlockData return all blocks of blockchain.
	GetAllBlockData() []string
	// Close closes blockchain.
	Close() error
}

type MockupExplorerServer struct {
	config     config.Config
	blockchain Blockchain
}

func NewMockupExplorerServer(config config.Config, blockchain Blockchain) ExplorerServer {
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
