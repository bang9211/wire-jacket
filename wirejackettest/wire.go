//go:build wireinject
// +build wireinject

package wirejackettest

import (
	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/bang9211/wire-jacket/wirejackettest/defaultexplorerserver"
	"github.com/bang9211/wire-jacket/wirejackettest/defaultrestapiserver"
	"github.com/bang9211/wire-jacket/wirejackettest/ossiconesblockchain"
	"github.com/google/wire"
)

//
// Dependency Injection List
//
// injectors stores module_name(key) with injector_func(value) using map.
// For wiring, name of implement using in config with injector function.
//
// Examples :
//
//	var injectors = map[string]interface{}{
// 		"ossiconesblockchain":	InjectOssiconesBlockchain,
// 	}
//
// 	var eagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// 	}
//
var injectors = map[string]interface{}{
	"ossiconesblockchain": InjectOssiconesBlockchain,
}

var eagerInjectors = map[string]interface{}{
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
func InjectOssiconesBlockchain(config config.Config) (ossiconesblockchain.Blockchain, error) {
	wire.Build(ossiconesblockchain.GetOrCreate)
	return nil, nil
}

// InjectDefaultExplorerServer injects dependencies and inits of ExplorerServer.
func InjectDefaultExplorerServer(
	config config.Config,
	blockchain ossiconesblockchain.Blockchain,
) (defaultexplorerserver.ExplorerServer, error) {
	wire.Build(defaultexplorerserver.GetOrCreate)
	return nil, nil
}

// InjectDefaultRESTAPIServer injects dependencies and inits of APiServer.
func InjectDefaultRESTAPIServer(
	config config.Config,
	blockchain ossiconesblockchain.Blockchain,
) (defaultrestapiserver.RESTAPIServer, error) {
	wire.Build(defaultrestapiserver.GetOrCreateDefault)
	return nil, nil
}
