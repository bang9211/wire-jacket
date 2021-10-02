package wirejacket

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/bang9211/wire-jacket/internal/utils"
)

// All the module should have Close().
type Module interface {
	// Close closes module gracefully.
	Close() error
}

// WireJacket structure. the jacket of the wires(injectors).
type WireJacket struct {
	config                 config.Config
	injectors              map[string]interface{}
	eagerInjectors         map[string]interface{}
	modules                map[string]Module
	sortedModulesByCreated []Module
	activatingModuleNames  []string
}

// New creates empty WireJacket.
// If you want to use more than one WireJacket on the same system,
// Use NewWithServiceName with unique serviceName instead of New().
// By default, WireJacket reads list of module name to activate
// in 'modules' value of config.
//
// But Wire-Jacket considered The Twelve Factors. Config can be
// overrided by envrionment variable.(see viperconfig.go)
// So, when using more than one WireJacket on the same system,
// each WireJacket should have a unique service name to avoid
// conflicting value of 'modules'.
//
// If serviceName no exists, WireJacket reads value of
// 'modules'.
// By default, WireJacket reads app.conf. Or you can specify
// file with '--config' flag.(see viperconfig.go)
//
// Example in app.conf without serviceName
//
// modules=mockup_database mockup_blockchain mockup_explorerserver mockup_restapiserver
//
// SetActivatingModules can specify activating modules without config.
// But this way is needed re-compile for changing module.
// The list of activating modules is used as key of injectors
// to call.
func New() *WireJacket {
	viperConfig := config.NewViperConfig()
	wj := &WireJacket{
		config:                 viperConfig,
		injectors:              map[string]interface{}{},
		eagerInjectors:         map[string]interface{}{},
		modules:                map[string]Module{"viperconfig": viperConfig},
		sortedModulesByCreated: []Module{viperConfig},
	}
	wj.activatingModuleNames = wj.readActivatingModules("")
	wj.activatingModuleNames = append(wj.activatingModuleNames, "viperconfig")

	return wj
}

// NewWithServiceName creates empty WireJacket.
// Make sure the serviceName is unique in same system.
// By default, WireJacket reads list of module name to activate
// in 'modules' value of config.
//
// But Wire-Jacket considered The Twelve Factors. Config can be
// overrided by envrionment variable.(see viperconfig.go)
// So, when using more than one WireJacket on the same system,
// each WireJacket should have a unique service name to avoid
// conflicting value of 'modules'.
//
// If serviceName exists, WireJacket reads value of
// '{serviceName}_modules' in config.
// By default, WireJacket reads app.conf. Or you can specify
// file with '--config' flag.(see viperconfig.go)
//
// Example in app.conf with serviceName(ossicones)
//
// ossicones_modules=mockup_database mockup_blockchain mockup_explorerserver mockup_restapiserver
//
// SetActivatingModules can specify activating modules without config.
// But this way is needed re-compile for changing module.
// The list of activating modules is used as key of injectors
// to call.
func NewWithServiceName(serviceName string) *WireJacket {
	viperConfig := config.NewViperConfig()
	wj := &WireJacket{
		config:                 viperConfig,
		injectors:              map[string]interface{}{},
		eagerInjectors:         map[string]interface{}{},
		modules:                map[string]Module{"viperconfig": viperConfig},
		sortedModulesByCreated: []Module{viperConfig},
	}
	wj.activatingModuleNames = wj.readActivatingModules(serviceName)
	wj.activatingModuleNames = append(wj.activatingModuleNames, "viperconfig")

	return wj
}

func (wj *WireJacket) readActivatingModules(serviceName string) []string {
	var activatingModuleNames []string
	if serviceName == "" {
		activatingModuleNames = wj.config.GetStringSlice(
			"modules", []string{},
		)
	} else {
		activatingModuleNames = wj.config.GetStringSlice(
			strings.ToLower(serviceName)+"_modules", []string{},
		)
	}

	return activatingModuleNames
}

// SetActivatingModules sets list of module's name.
// module's name is used as key of injector maps.
// It overwrites list of modules to activate.
func (wj *WireJacket) SetActivatingModules(moduleNames []string) {
	wj.activatingModuleNames = moduleNames
	wj.activatingModuleNames = append(wj.activatingModuleNames, "viperconfig")
}

// SetInjectors sets injectors to inject lazily.
// WireJacket maps module_name to injector as a key-value pairs.
// WireJacket tries to find module_name.
// if serviceName no exists, value of 'modules' in config.
// if serviceName exists, value of '{serviceName}_modules' in config.
//
//
// Example of app.conf (serviceName=ossicones) :
//
// ossicones_modules=mockup_database mockup_blockchain mockup_explorerserver mockup_restapiserver
//
//
// definition in wire.go
//
// func InjectViperConfig() (config.Config, error) { ... }
// func InjectOssiconesBlockchain(config config.Config) (blockchain.Blockchain, error) { ... }
// func InjectDefaultExplorerServer(config config.Config, blockchain blockchain.Blockchain) (explorerserver.ExplorerServer, error) { ...}
// func InjectDefaultRESTAPIServer(config config.Config, blockchain blockchain.Blockchain) (restapiserver.RESTAPIServer, error) { ...}
//
//
// injectors can be like this.
//
// var injectors = map[string]interface{}{
// 		"viperconfig":         InjectViperConfig,
// 		"ossiconesblockchain": InjectOssiconesBlockchain,
// }
//
// eagerInjectors can be like this.
//
// var eagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// }
//
//
// injectors will be injected lazily.
func (wj *WireJacket) SetInjectors(injectors map[string]interface{}) *WireJacket {
	wj.injectors = injectors
	return wj
}

// SetEagerInjectors sets injectors to inject eagerly.
// WireJacket maps module_name to injector as a key-value pairs.
// WireJacket tries to find module_name.
// if serviceName no exists, value of 'modules' in config.
// if serviceName exists, value of '{serviceName}_modules' in config.
//
//
// Example of app.conf (serviceName=ossicones) :
//
// ossicones_modules=mockup_database mockup_blockchain mockup_explorerserver mockup_restapiserver
//
//
// definition in wire.go
//
// func InjectViperConfig() (config.Config, error) { ... }
// func InjectOssiconesBlockchain(config config.Config) (blockchain.Blockchain, error) { ... }
// func InjectDefaultExplorerServer(config config.Config, blockchain blockchain.Blockchain) (explorerserver.ExplorerServer, error) { ...}
// func InjectDefaultRESTAPIServer(config config.Config, blockchain blockchain.Blockchain) (restapiserver.RESTAPIServer, error) { ...}
//
//
// injectors can be like this.
//
// var injectors = map[string]interface{}{
// 		"viperconfig":         InjectViperConfig,
// 		"ossiconesblockchain": InjectOssiconesBlockchain,
// }
//
// eagerInjectors can be like this.
//
// var eagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// }
//
//
// injectors will be injected eagerly.
func (wj *WireJacket) SetEagerInjectors(injectors map[string]interface{}) *WireJacket {
	wj.eagerInjectors = injectors
	return wj
}

// getInjector gets injector in eagerInjectors or injectors
func (wj *WireJacket) getInjector(moduleName string) interface{} {
	injector := wj.eagerInjectors[moduleName]
	if injector == nil {
		injector = wj.injectors[moduleName]
		if injector == nil {
			return nil
		}
	}
	return injector
}

// getInjector gets all injectors in eagerInjectors and injectors
func (wj *WireJacket) getInjectors() map[string]interface{} {
	injectors := map[string]interface{}{}
	for moduleName, injector := range wj.eagerInjectors {
		injectors[moduleName] = injector
	}
	for moduleName, injector := range wj.injectors {
		injectors[moduleName] = injector
	}

	return injectors
}

// AddInjector adds injector function to the lazy injection list.
func (wj *WireJacket) AddInjector(moduleName string, injector interface{}) {
	if reflect.TypeOf(injector).Kind() == reflect.Func {
		wj.injectors[moduleName] = injector
	}
}

// AddInjector adds injector function to the eager injection list.
func (wj *WireJacket) AddEagerInjector(moduleName string, injector interface{}) {
	if reflect.TypeOf(injector).Kind() == reflect.Func {
		wj.eagerInjectors[moduleName] = injector
	}
}

// DoWire does wiring of wires(injectors).
// It calls eagerInjectors as finding(if no exists, loading) and injecting dependencies.
func (wj *WireJacket) DoWire() error {
	if len(wj.getInjectors()) == 0 {
		return fmt.Errorf("no injectors to wire")
	}
	if len(wj.activatingModuleNames) == 1 { //default viperconfig
		return fmt.Errorf("no activating modules to wire")
	}
	for moduleName, eagerInjector := range wj.eagerInjectors {
		err := wj.loadModule(moduleName, eagerInjector)
		if err != nil {
			return fmt.Errorf("[%s] %s", moduleName, err)
		}
	}

	return nil
}

func (wj *WireJacket) loadModule(moduleName string, injector interface{}) error {
	//already exists
	if wj.modules[moduleName] != nil {
		return nil
	}
	if !utils.IsContain(wj.activatingModuleNames, moduleName) {
		return fmt.Errorf("no activating module name for injector(%s), in %s",
			moduleName,
			wj.activatingModuleNames)
	}

	// get dependencies
	injectorFunc := reflect.ValueOf(injector)
	dependencies, err := wj.getDependencies(moduleName, injectorFunc.Type())
	if err != nil {
		return err
	}

	// call injector
	returnVal := injectorFunc.Call(dependencies)
	module, err := wj.checkInjectionResult(returnVal)
	if err != nil {
		return err
	}

	// set module
	wj.modules[moduleName] = module
	pushModule(&wj.sortedModulesByCreated, module)

	return nil
}

func (wj *WireJacket) getDependencies(
	moduleName string,
	injectorFuncType reflect.Type) ([]reflect.Value, error) {
	dependencyTypeList := wj.getDependencyTypeList(injectorFuncType)

	dependencies, err := wj.loadAndGetDependencies(moduleName, dependencyTypeList)
	if err != nil {
		return nil, err
	}

	return dependencies, nil
}

func (wj *WireJacket) loadAndGetDependencies(
	moduleName string,
	dependencyTypeList []reflect.Type) ([]reflect.Value, error) {
	dependencies := []reflect.Value{}
	for _, dependencyType := range dependencyTypeList {
		dependencyPtr := wj.findDependency(dependencyType)
		if dependencyPtr != nil {
			dependencies = append(dependencies, *dependencyPtr)
		} else {
			// find injector to create dependency (return type check)
			moduleName, injector := wj.findInjector(dependencyType)
			if injector == nil {
				return nil, fmt.Errorf("failed to find injector of dependency(%s)", moduleName)
			}

			// load dependency using injector
			err := wj.loadModule(moduleName, injector)
			if err != nil {
				return nil, fmt.Errorf("failed to load module of dependency(%s)", moduleName)
			}

			// add loaded dependency by injector
			dependencies = append(dependencies, reflect.ValueOf(wj.modules[moduleName]))
		}
	}

	return dependencies, nil
}

func (wj *WireJacket) findInjector(dependencyType reflect.Type) (string, interface{}) {
	injectors := wj.getInjectors()
	for moduleName, injector := range injectors {
		injectorFunc := reflect.ValueOf(injector)
		injectorFuncType := injectorFunc.Type()
		if injectorFuncType.NumOut() > 0 &&
			injectorFuncType.Out(0).Name() == dependencyType.Name() &&
			injectorFuncType.Out(0).PkgPath() == dependencyType.PkgPath() {
			return moduleName, injector
		}
	}

	return "", nil
}

func (wj *WireJacket) getDependencyTypeList(injectorFuncType reflect.Type) []reflect.Type {
	typeList := []reflect.Type{}
	for i := 0; i < injectorFuncType.NumIn(); i++ {
		dependency := injectorFuncType.In(i)
		typeList = append(typeList, dependency)
	}
	return typeList
}

func (wj *WireJacket) findDependency(dependencyType reflect.Type) *reflect.Value {
	for moduleName, module := range wj.modules {
		if utils.IsContain(wj.activatingModuleNames, moduleName) {
			moduleValue := reflect.ValueOf(module)
			if moduleValue.CanConvert(dependencyType) {
				return &moduleValue
			}
		}
	}
	return nil
}

func (wj *WireJacket) checkInjectionResult(returnVal []reflect.Value) (Module, error) {
	if len(returnVal) != 1 && len(returnVal) != 2 {
		return nil, fmt.Errorf(
			"invalid inject function format len(return) : %d", len(returnVal))
	}
	var module Module
	var ok bool
	if len(returnVal) == 1 { // return (module)
		if !returnVal[0].IsValid() || !returnVal[0].CanInterface() {
			return nil, fmt.Errorf(
				"returnVal(%s) can't be interface",
				returnVal[0],
			)
		}
		module, ok = returnVal[0].Interface().(Module)
		if !ok {
			return nil, fmt.Errorf(
				"failed to cast returnVal(%s) to Module", returnVal[0])
		}
	} else { // return (module, error)
		if !returnVal[1].IsValid() || !returnVal[1].CanInterface() {
			return nil, fmt.Errorf(
				"failed to cast error(%s) to interface", returnVal[1])
		}
		err := returnVal[1].Interface()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to inject : %s", err)
		}
		if !returnVal[0].IsValid() || !returnVal[0].CanInterface() {
			return nil, fmt.Errorf(
				"failed to cast returnVal(%s) to interface", returnVal[0])
		}
		module, ok = returnVal[0].Interface().(Module)
		if !ok {
			return nil, fmt.Errorf(
				"failed to cast returnVal(%s) to Module", returnVal[0])
		}
	}
	return module, nil
}

// GetConfig returns config object.
func (wj *WireJacket) GetConfig() config.Config {
	return wj.config
}

// GetModule finds module using moduleName and returns module if exists.
// If no exists, it tries to create module using injector and returns.
func (wj *WireJacket) GetModule(moduleName string) interface{} {
	module := wj.modules[moduleName]
	if module != nil {
		return module
	}
	injector := wj.getInjector(moduleName)
	if injector == nil {
		return nil
	}
	wj.loadModule(moduleName, injector)

	return wj.modules[moduleName]
}

// GetModuleByType finds module using interfaceType(pointer of interface)
// and returns module if exists.
//
// Example (process_name=ossicones) :
//
// config := wj.GetModuleByType((*config.Config)(nil))
//
// If no exists, it tries to create module using injector and returns.
// This may return undesirable results if there are other implementations
// that may use the same interface. Use only if you are sure that there
// are no overlapping interfaces.
func (wj *WireJacket) GetModuleByType(interfaceType interface{}) interface{} {
	if interfaceType == nil {
		return nil
	}
	moduleType := reflect.TypeOf(interfaceType).Elem()
	moduleName, injector := wj.findInjector(moduleType)
	wj.loadModule(moduleName, injector)
	for _, module := range wj.modules {
		moduleValue := reflect.ValueOf(module)
		if moduleValue.CanConvert(moduleType) {
			return module
		}
	}

	return nil
}

// Close closes all the modules gracefully
func (wj *WireJacket) Close() error {
	moduleLen := len(wj.sortedModulesByCreated)
	for i := 0; i < moduleLen; i++ {
		module := popModule(&wj.sortedModulesByCreated)
		err := module.Close()
		if err != nil {
			log.Printf("failed to close module(%s) : %s", reflect.ValueOf(module).Type(), err)
		}
	}

	return nil
}

func pushModule(modules *[]Module, module Module) {
	*modules = append(*modules, module)
}

func popModule(modules *[]Module) Module {
	ret := (*modules)[0]
	if len(*modules) > 0 {
		*modules = (*modules)[1:]
	}
	return ret
}
