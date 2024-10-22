//go:build wireinject
// +build wireinject

package mockup

import (
	viperjacket "github.com/bang9211/viper-jacket"
	"github.com/google/wire"
)

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

//
// Dependency wiring should be specify in wire.go.
//
// Inject functions can have several dependency parameters
// and should have two returns(interface, error).
// Only structure type is allowed, non-structure(int, string, ...) is not allowed for injection.
//
// Function Form :
//
// - func Inject{Implement}() {Interface} {}
// - func Inject{Implement}() ({Interface}, error) {}
// - func Inject{Implement}({Interface}) {Interface} {}
// - func Inject{Implement}({Interface}) ({Interface}, error) {}
//
// Examples :
//
// - func InjectViperJacket() viperjacket.Config {}
// - func InjectViperJacket() (viperjacket.Config, error) {}
// - func InjectOssiconesBlockChain(config viperjacket.Config) blockchain.Blockchain {}
// - func InjectOssiconesBlockChain(config viperjacket.Config) (blockchain.Blockchain, error) {}
//

// InjectMockupDB injects dependencies and inits of Database.
func InjectMockupDB(config viperjacket.Config) (Database, error) {
	wire.Build(NewMockupDB)
	return nil, nil
}

// InjectMockupBlockchain injects dependencies and inits of Blockchain.
func InjectMockupBlockchain(db Database) (Blockchain, error) {
	wire.Build(NewMockupBlockchain)
	return nil, nil
}

// InjectMockupExplorerServer injects dependencies and inits of ExplorerServer.
func InjectMockupExplorerServer(
	config viperjacket.Config,
	blockchain Blockchain,
) (ExplorerServer, error) {
	wire.Build(NewMockupExplorerServer)
	return nil, nil
}

// InjectMockupRESTAPIServer injects dependencies and inits of RESTAPIServer.
func InjectMockupRESTAPIServer(
	config viperjacket.Config,
	blockchain Blockchain,
) (RESTAPIServer, error) {
	wire.Build(NewMockupRESTAPIServer)
	return nil, nil
}

// InjectMockupInvalidReturnTest injects dependencies and inits of
// RESTAPIServer and return invalid format.
func InjectMockupInvalidReturnTest(
	config viperjacket.Config,
	blockchain Blockchain,
) (RESTAPIServer, func(), error) {
	wire.Build(NewMockupRESTAPIServer)
	return nil, func() {}, nil
}

// InjectMockupInvalidImplTest injects test dependency for test.
func InjectMockupInvalidImplTest() (TestInterface, error) {
	wire.Build(NewTestImplement)
	return nil, nil
}
