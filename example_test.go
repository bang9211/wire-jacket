package wirejacket_test

import (
	"fmt"
	"log"

	wirejacket "github.com/bang9211/wire-jacket"
	"github.com/bang9211/wire-jacket/internal/config"
)

// ==============================================
// Database Interface - MockupDB Implment example
// ==============================================
type Database interface {
	// Connect DB.
	Connect() error
	// Close closes the REST API Server.
	Close() error
}

type MockupDB struct {
	config config.Config
}

func NewMockupDB(config config.Config) Database {
	return &MockupDB{config: config}
}

func (mdb *MockupDB) Connect() error {
	log.Printf("connect : %s", mdb.config.GetString("address", "localhost:3306"))
	return nil
}

func (mdb *MockupDB) Close() error {
	// drs = nil
	return nil
}

// =========================================================
// Blockchain Interface - MockupBlockchain Implement example
// =========================================================
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

// =======================================
// wire_gen.go(Wire generated code) example
// =======================================

// InjectMockupDB injects dependencies and inits of Database.
func InjectMockupDB(config2 config.Config) (Database, error) {
	database := NewMockupDB(config2)
	return database, nil
}

// InjectMockupBlockchain injects dependencies and inits of Blockchain.
func InjectMockupBlockchain(db Database) (Blockchain, error) {
	blockchain := NewMockupBlockchain(db)
	return blockchain, nil
}

var Injectors = map[string]interface{}{
	"mockup_database": InjectMockupDB,
}

var EagerInjectors = map[string]interface{}{
	"mockup_blockchain": InjectMockupBlockchain,
}

// ==================
// User code examples
// ==================

// Default use case to use New().
// Wire-Jacket defaultly uses 'app.conf' for setting modules
// to activate. Or you can use the flag '--config {file_name}'.
func Example_New() {
	// Create wirejacket and set injectors.
	wj := wirejacket.New().
		SetEagerInjectors(EagerInjectors).
		SetInjectors(Injectors)

	// Inject eager injectors.
	if err := wj.DoWire(); err != nil {
		log.Fatal(err)
	}

	// Check value of modules in app.conf.
	// Or You can set modules using SetActivatingModules() directly.
	wj.SetActivatingModules([]string{"mockup_blockchain", "mockup_databse"})

	// Get module and use.
	blockchain := wj.GetModule("mockup_blockchain").(Blockchain)
	fmt.Println(blockchain.GetBlocks())
}

// Second use case to use NewWithServiceName().
func Example_NewWithServiceName() {
	// Create wirejacket with serviceName.
	wj := wirejacket.NewWithServiceName("example_service")

	// You can also add injector directly, instead of SetInjectors().
	wj.AddInjector("mockup_database", InjectMockupDB)
	wj.AddEagerInjector("mockup_blockchain", InjectMockupBlockchain)

	// Check value of modules in app.conf.
	// Or You can set modules using SetActivatingModules() directly.
	wj.SetActivatingModules([]string{"mockup_blockchain", "mockup_databse"})

	// You can also get module without DoWire().
	// the dependencies of the module are injected automatically.
	blockchain := wj.GetModule("mockup_blockchain").(Blockchain)
	fmt.Println(blockchain.GetBlocks())
}
