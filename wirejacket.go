package wirejacket

import (
	"fmt"
	"log"
	"os"
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
	serviceName            string
	config                 config.Config
	injectors              map[string]interface{}
	eagerInjectors         map[string]interface{}
	modules                map[string]Module
	sortedModulesByCreated []Module
	activatingModuleNames  []string
}

// New creates empty WireJacket. serviceName uses as prefix of config.
// Wirejacket's config can be overrided by envrionment variable.
// So it needs unique serviceName to avoid collision.
// By default, it uses {serviceName}.conf file for reading list of 
// activating module.
//
// Example in {serviceName}.conf
//
// {serviceName}_activating_modules=ossiconesblockchain viperconfig defaultexplorerserver defaultrestapiserver
// 
// If {serviceName}.conf file no exists, SetActivatingModules must be 
// called to specify activating modules.
// The list of activating module is used as key of injectors to call.
func New(serviceName string) (*WireJacket, error) {
	wj := &WireJacket{
		serviceName:            serviceName,
		config:                 config.NewViperConfig(),
		injectors:              map[string]interface{}{},
		eagerInjectors:         map[string]interface{}{},
		modules:                map[string]Module{},
		sortedModulesByCreated: []Module{},
	}
	wj.activatingModuleNames = readActivatingModules(wj.config)

	return wj, nil
}

// NewWithInjectors creates WireJacket with injectors.
// Wirejacket's config can be overrided by envrionment variable.
// So it needs unique serviceName to avoid collision.
// By default, it uses {serviceName}.conf file for reading list 
// of activating module.
//
// Example in {serviceName}.conf
//
// {serviceName}_activating_modules=ossiconesblockchain viperconfig defaultexplorerserver defaultrestapiserver
// 
// If {serviceName}.conf file no exists, SetActivatingModules must be 
// called to specify activating modules.
// The list of activating module is used as key of injectors to call.
func NewWithInjectors(
	serviceName string,
	injectors map[string]interface{},
	eagerInjectors map[string]interface{}) (*WireJacket, error) {
	wj := &WireJacket{
		serviceName:            serviceName,
		injectors:              injectors,
		eagerInjectors:         eagerInjectors,
		modules:                map[string]Module{},
		config:                 config.NewViperConfig(),
		sortedModulesByCreated: []Module{},
	}
	wj.activatingModuleNames = readActivatingModules(wj.config)

	return wj, nil
}

func readActivatingModules(config config.Config) []string {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	activatingModuleNames := config.GetStringSlice(
		strings.ToLower(os.Args[0])+"_activating_modules",
		[]string{},
		// defaultActivatingModules[:], // array to slice
	)

	return activatingModuleNames
}

// SetActivatingModules sets case-insensitive list of module's name.
// module's name is used as key of injector maps.
// It can be overrided by value of {serviceName}_activating_modules in {serviceName}.conf file.
func (w *WireJacket) SetActivatingModules(moduleNames []string) {
	w.activatingModuleNames = moduleNames
}

// SetInjectors set injectors to inject lazily.
// Wire has two basic concepts: providers and injectors.
// WireJacket maps injector to module_name and injector as a key-value pairs.
// The module_name can be found in the config file.
// By default, WireJacket tries to find module_name in
// {process_name}_activating_modules of {process_name}.conf file.
//
//
// Example (process_name=ossicones) :
//
//
// # declaration in ossicones.conf
// ossicones_activating_modules=ossiconesblockchain viperconfig defaultexplorerserver defaultrestapiserver
//
//
// # definition in wire.go
// func InjectViperConfig() (config.Config, error) { ... }
// func InjectOssiconesBlockchain(config config.Config) (blockchain.Blockchain, error) { ... }
// func InjectDefaultExplorerServer(config config.Config, blockchain blockchain.Blockchain) (explorerserver.ExplorerServer, error) { ...}
// func InjectDefaultRESTAPIServer(config config.Config, blockchain blockchain.Blockchain) (restapiserver.RESTAPIServer, error) { ...}
//
//
// # injectors can be like this.
// var injectors = map[string]interface{}{
// 		"viperconfig":         InjectViperConfig,
// 		"ossiconesblockchain": InjectOssiconesBlockchain,
// }
//
// # eagerInjectors can be like this.
// var eagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// }
//
//
// injectors will be injected lazily.
func (wj *WireJacket) SetInjectors(injectors map[string]interface{}) {
	wj.injectors = injectors
}

// SetEagerInjectors set injectors to inject eagerly.
// Wire has two basic concepts: providers and injectors.
// WireJacket maps injector to module_name and injector as a key-value pairs.
// The module_name can be found in the config file.
// By default, WireJacket tries to find module_name in
// {process_name}_activating_modules of {process_name}.conf file.
//
//
// Example (process_name=ossicones) :
//
//
// # declaration in ossicones.conf
// ossicones_activating_modules=ossiconesblockchain viperconfig defaultexplorerserver defaultrestapiserver
//
//
// # definition in wire.go
// func InjectViperConfig() (config.Config, error) { ... }
// func InjectOssiconesBlockchain(config config.Config) (blockchain.Blockchain, error) { ... }
// func InjectDefaultExplorerServer(config config.Config, blockchain blockchain.Blockchain) (explorerserver.ExplorerServer, error) { ...}
// func InjectDefaultRESTAPIServer(config config.Config, blockchain blockchain.Blockchain) (restapiserver.RESTAPIServer, error) { ...}
//
//
// # injectors can be like this.
// var injectors = map[string]interface{}{
// 		"viperconfig":         InjectViperConfig,
// 		"ossiconesblockchain": InjectOssiconesBlockchain,
// }
//
// # eagerInjectors can be like this.
// var eagerInjectors = map[string]interface{}{
// 		"defaultexplorerserver": InjectDefaultExplorerServer,
// 		"defaultrestapiserver":  InjectDefaultRESTAPIServer,
// }
//
//
// injectors will be injected eagerly.
func (wj *WireJacket) SetEagerInjectors(injectors map[string]interface{}) {
	wj.eagerInjectors = injectors
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

// AddInjector add injector function to the lazy injection list.
func (wj *WireJacket) AddInjector(moduleName string, injector interface{}) {
	if reflect.TypeOf(injector).Kind() == reflect.Func {
		wj.injectors[moduleName] = injector
	}
}

// AddInjector add injector function to the eager injection list.
func (wj *WireJacket) AddEagerInjector(moduleName string, injector interface{}) {
	if reflect.TypeOf(injector).Kind() == reflect.Func {
		wj.eagerInjectors[moduleName] = injector
	}
}

// DoWire does wiring of wires(injectors).
// It calls eagerInjectors as finding(if no exists, loading) and injecting dependencies.
func (wj *WireJacket) DoWire() error {
	for moduleName, eagerInjector := range wj.eagerInjectors {
		if utils.IsContain(wj.activatingModuleNames, moduleName) {
			err := wj.loadModule(moduleName, eagerInjector)
			if err != nil {
				return fmt.Errorf("[%s] %s", moduleName, err)
			}
		}
	}

	// if eagerInjectors no exists, wire all.
	if len(wj.eagerInjectors) == 0 {
		wj.loadAllModules()
	}

	return nil
}

func (wj *WireJacket) loadModule(moduleName string, injector interface{}) error {
	//already exists
	if wj.modules[moduleName] != nil {
		return nil
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

func (wj *WireJacket) getDependencies(moduleName string, injectorFuncType reflect.Type) ([]reflect.Value, error) {
	dependencyTypeList := wj.getDependencyTypeList(injectorFuncType)

	dependencies := wj.findDependencies(dependencyTypeList)
	if dependencies == nil {
		var err error
		dependencies, err = wj.loadAndGetDependencies(moduleName, dependencyTypeList)
		if err != nil {
			return nil, err
		}
	}
	return dependencies, nil
}

func (wj *WireJacket) findDependencies(dependencyTypeList []reflect.Type) []reflect.Value {
	dependencies := []reflect.Value{}
	for i := 0; i < len(dependencyTypeList); i++ {
		dependencyType := dependencyTypeList[i]
		moduleValue := wj.findModuleValue(dependencyType)
		if moduleValue == nil {
			return nil
		}
		dependencies = append(dependencies, *moduleValue)
	}
	return dependencies
}

func (wj *WireJacket) loadAndGetDependencies(moduleName string, dependencyTypeList []reflect.Type) ([]reflect.Value, error) {
	var err error
	dependencies := []reflect.Value{}

	for _, dependencyType := range dependencyTypeList {
		// find injector of dependency in injectors (return type check)
		moduleName, injector := wj.findInjector(dependencyType)
		if injector == nil {
			return nil, fmt.Errorf("failed to find injector of dependency(%s)", moduleName)
		}

		// load module using injector
		err = wj.loadModule(moduleName, injector)
		if err != nil {
			return nil, fmt.Errorf("failed to load module of dependency(%s)", moduleName)
		}

		// get module as dependency
		dependencies = append(dependencies, reflect.ValueOf(wj.modules[moduleName]))
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

func (wj *WireJacket) findModuleValue(dependencyType reflect.Type) *reflect.Value {
	for _, module := range wj.modules {
		moduleValue := reflect.ValueOf(module)
		if moduleValue.CanConvert(dependencyType) {
			return &moduleValue
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
		if !returnVal[0].CanInterface() {
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
		if !returnVal[1].CanInterface() {
			return nil, fmt.Errorf(
				"failed to cast error(%s) to interface", returnVal[1])
		}
		err := returnVal[1].Interface()
		if err != nil {
			return nil, fmt.Errorf(
				"failed to inject : %s", err)
		}
		if !returnVal[0].CanInterface() {
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

func (wj *WireJacket) loadAllModules() error {
	NotActivatedList := make([]string, len(wj.activatingModuleNames))
	copy(NotActivatedList, wj.activatingModuleNames)
	activatedList := []string{}
	tryCount := 0

	for len(NotActivatedList) > 0 && tryCount < len(NotActivatedList)*len(NotActivatedList) {
		for _, moduleName := range NotActivatedList {
			injector := wj.getInjector(moduleName)
			if injector == nil {
				return fmt.Errorf("failed to get injector of %s", moduleName)
			}
			injectorFunc := reflect.ValueOf(injector)
			injectorFuncType := injectorFunc.Type()
			dependencies, err := wj.getDependencies(moduleName, injectorFuncType)
			if err != nil {
				return err
			}
			if dependencies != nil {
				returnVal := injectorFunc.Call(dependencies)
				module, err := wj.checkInjectionResult(returnVal)
				if err != nil {
					return fmt.Errorf("[%s] %s", moduleName, err)
				}
				wj.modules[moduleName] = module
				pushModule(&wj.sortedModulesByCreated, module)
				activatedList = append(activatedList, moduleName)
			}
		}
		for _, activated := range activatedList {
			NotActivatedList = utils.RemoveElement(NotActivatedList, activated)
		}
		tryCount++
	}

	return nil
}

// GetConfig returns config object.
func (wj *WireJacket) GetConfig() config.Config {
	return wj.config
}

// GetModule finds module using moduleName and returns module if exists.
// If no exists, it tries to create module using injector and returns.
func (wj *WireJacket) GetModule(moduleName string) interface{} {
	injector := wj.getInjector(moduleName)
	if injector == nil {
		return nil
	}
	wj.loadModule(moduleName, injector)

	return wj.modules[moduleName]
}

// GetModuleByType finds module using interfaceType and returns module if exists.
// If no exists, it tries to create module using injector and returns.
func (wj *WireJacket) GetModuleByType(interfaceType interface{}) interface{} {
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
