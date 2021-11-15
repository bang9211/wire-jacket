// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package mockup

import (
	"github.com/bang9211/viper-jacket"
)

// Injectors from wire.go:

// InjectMockupDB injects dependencies and inits of Database.
func InjectMockupDB(config viperjacket.Config) (Database, error) {
	database := NewMockupDB(config)
	return database, nil
}

// InjectMockupBlockchain injects dependencies and inits of Blockchain.
func InjectMockupBlockchain(db Database) (Blockchain, error) {
	blockchain := NewMockupBlockchain(db)
	return blockchain, nil
}

// InjectMockupExplorerServer injects dependencies and inits of ExplorerServer.
func InjectMockupExplorerServer(config viperjacket.Config, blockchain Blockchain) (ExplorerServer, error) {
	explorerServer := NewMockupExplorerServer(config, blockchain)
	return explorerServer, nil
}

// InjectMockupRESTAPIServer injects dependencies and inits of RESTAPIServer.
func InjectMockupRESTAPIServer(config viperjacket.Config, blockchain Blockchain) (RESTAPIServer, error) {
	restapiServer := NewMockupRESTAPIServer(config, blockchain)
	return restapiServer, nil
}

// InjectMockupInvalidReturnTest injects dependencies and inits of
// RESTAPIServer and return invalid format.
func InjectMockupInvalidReturnTest(config viperjacket.Config, blockchain Blockchain) (RESTAPIServer, func(), error) {
	restapiServer := NewMockupRESTAPIServer(config, blockchain)
	return restapiServer, func() {
	}, nil
}

// InjectMockupInvalidImplTest injects test dependency for test.
func InjectMockupInvalidImplTest() (TestInterface, error) {
	testInterface := NewTestImplement()
	return testInterface, nil
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
