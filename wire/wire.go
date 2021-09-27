//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/bang9211/wire-jacket/internal/config"
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
// 		"ossiconesblockchain":	InjectOssiconesBlockchain,
// 	}
//
// 	var EagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// 	}
//
var Injectors = map[string]interface{}{
	"ossiconesblockchain": InjectOssiconesBlockchain,
}

var EagerInjectors = map[string]interface{}{
	"defaultexplorerserver": InjectDefaultExplorerServer,
	"defaultrestapiserver":  InjectDefaultRESTAPIServer,
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
// - func InjectViperConfig() config.Config {}
// - func InjectViperConfig() (config.Config, error) {}
// - func InjectOssiconesBlockChain(config config.Config) blockchain.Blockchain {}
// - func InjectOssiconesBlockChain(config config.Config) (blockchain.Blockchain, error) {}
//

// InjectOssiconesBlockchain injects dependencies and inits of Blockchain.
func InjectOssiconesBlockchain(config config.Config) (Blockchain, error) {
	wire.Build(GetOrCreateOssiconesBlockchain)
	return nil, nil
}

// InjectDefaultExplorerServer injects dependencies and inits of ExplorerServer.
func InjectDefaultExplorerServer(
	config config.Config,
	blockchain Blockchain,
) (ExplorerServer, error) {
	wire.Build(GetOrCreateDefaultExplorerServer)
	return nil, nil
}

// InjectDefaultRESTAPIServer injects dependencies and inits of APiServer.
func InjectDefaultRESTAPIServer(
	config config.Config,
	blockchain Blockchain,
) (RESTAPIServer, error) {
	wire.Build(GetOrCreateDefaultRESTAPIServer)
	return nil, nil
}
