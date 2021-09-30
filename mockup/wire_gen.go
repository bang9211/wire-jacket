// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package mockup

import (
	"github.com/bang9211/wire-jacket/internal/config"
)

// Injectors from wire.go:

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

// InjectMockupExplorerServer injects dependencies and inits of ExplorerServer.
func InjectMockupExplorerServer(config2 config.Config, blockchain Blockchain) (ExplorerServer, error) {
	explorerServer := NewMockupExplorerServer(config2, blockchain)
	return explorerServer, nil
}

// InjectMockupRESTAPIServer injects dependencies and inits of RESTAPIServer.
func InjectMockupRESTAPIServer(config2 config.Config, blockchain Blockchain) (RESTAPIServer, error) {
	restapiServer := NewMockupRESTAPIServer(config2, blockchain)
	return restapiServer, nil
}

// wire.go:

//
// Dependency Injection List
//
// Injectors stores module_name(key) with injector_func(value) using map.
// For wiring, name of implement using in config with injector function.
//
// Examples :
//
//	var Injectors = map[string]interface{}{
// 		"mockup_database":   InjectMockupDB,
// 		"mockup_blockchain": InjectMockupBlockchain,
// 	}
//
// 	var EagerInjectors = map[string]interface{}{
//		"mockup_explorerserver": InjectMockupExplorerServer,
//		"mockup_restapiserver":  InjectMockupRESTAPIServer,
// 	}
//
var Injectors = map[string]interface{}{
	"mockup_database":   InjectMockupDB,
	"mockup_blockchain": InjectMockupBlockchain,
}

var EagerInjectors = map[string]interface{}{
	"mockup_explorerserver": InjectMockupExplorerServer,
	"mockup_restapiserver":  InjectMockupRESTAPIServer,
}
