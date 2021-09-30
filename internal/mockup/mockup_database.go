package mockup

import (
	"log"

	"github.com/bang9211/wire-jacket/internal/config"
)

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
