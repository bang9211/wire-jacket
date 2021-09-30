package mockup

type Block interface {
	GetData() string
}

type MockupBlock struct {
	data string
}

func (mb *MockupBlock) GetData() string {
	return mb.data
}

type Blockchain interface {
	// Init inits blockchain.
	Init() error
	// AddBlock adds data to blockchain.
	AddBlock(data string) error
	// GetBlocks gets all the blocks.
	GetBlocks() []Block
	// Close closes blockchain.
	Close() error
}

var genesisBlockData = "Genesis Block Data"

type MockupBlockchain struct {
	db     Database
	blocks []Block
}

func NewMockupBlockchain(db Database) Blockchain {
	return &MockupBlockchain{db: db, blocks: []Block{}}
}

func (mbc *MockupBlockchain) Init() error {
	mbc.db.Connect()
	mbc.AddBlock(genesisBlockData)
	return nil
}

func (mbc *MockupBlockchain) AddBlock(data string) error {
	mbc.blocks = append(mbc.blocks, &MockupBlock{data: data})
	return nil
}

func (mbc *MockupBlockchain) GetBlocks() []Block {
	return mbc.blocks
}

func (mbc *MockupBlockchain) Close() error {
	return nil
}
