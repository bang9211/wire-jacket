package wirejacket

import (
	"fmt"
	"log"
	"reflect"
)

// All the module should have Close().
type Module interface {
	Close() error
}

// WireJacket struct.
type WireJacket struct {
	config                Config
	injectors             map[string]interface{}
	modules               map[string]Module
	activatingModuleNames []string
}

// New creates empty WireJacket.
func New() (*WireJacket, error) {
	wj := &WireJacket{
		config:    NewViperConfig(),
		injectors: map[string]interface{}{},
		modules:   map[string]Module{},
	}

	return wj, nil
}

// NewWithInjectors creates WireJacket with injectors.
func NewWithInjectors(
	injectors map[string]interface{},
	eagerInjectors map[string]interface{}) (*WireJacket, error) {
	wj := &WireJacket{
		injectors: injectors,
		modules:   map[string]Module{},
		config:    NewViperConfig(),
	}
	wj.activatingModuleNames = readActivatingModules(wj.config)

	for moduleName, eagerInjector := range eagerInjectors {
		if IsContain(wj.activatingModuleNames, moduleName) {
			err := wj.loadModule(moduleName, eagerInjector)
			if err != nil {
				return nil, fmt.Errorf("[%s] %s", moduleName, err)
			}
		}
	}

	return wj, nil
}

func (wj *WireJacket) loadModule(moduleName string, injector interface{}) error {
	var err error

	method := reflect.ValueOf(injector)
	methodType := method.Type()
	dependencies, satisfied := wj.getDependencies(methodType)
	if !satisfied {
		dependencies, err = wj.loadDependencies(moduleName, methodType)
		if err != nil {
			return err
		}
	}

	returnVal := method.Call(dependencies)
	module, err := wj.checkInjectionResult(returnVal)
	if err != nil {
		return err
	}
	wj.modules[moduleName] = module

	return nil
}

func (wj *WireJacket) getDependencies(methodType reflect.Type) ([]reflect.Value, bool) {
	dependencies := []reflect.Value{}
	for i := 0; i < methodType.NumIn(); i++ {
		dependencyType := methodType.In(i)
		find := false
		for _, module := range wj.modules {
			moduleValue := reflect.ValueOf(module)
			if moduleValue.CanConvert(dependencyType) {
				dependencies = append(dependencies, moduleValue)
				find = true
				break
			}
		}
		if !find {
			return nil, false
		}
	}
	return dependencies, true
}

func (wj *WireJacket) loadDependencies(
	moduleName string,
	methodType reflect.Type) ([]reflect.Value, error) {
	var err error
	dependencies := []reflect.Value{}

	dependencyTypeList := wj.getParamTypeList(methodType)
	for _, dependencyType := range dependencyTypeList {
		// find injector of dependency in injectors (return type check)
		moduleName, injector := wj.findInjector(dependencyType)
		if injector == nil {
			return nil, fmt.Errorf("failed to find injector of dependency(%s)", moduleName)
		}

		// loadModule using injector
		err = wj.loadModule(moduleName, injector)
		if err != nil {
			return nil, fmt.Errorf("failed to load module of dependency(%s)", moduleName)
		}
		dependencies = append(dependencies, reflect.ValueOf(wj.modules[moduleName]))
	}

	return dependencies, nil
}

func (wj *WireJacket) findInjector(dependencyType reflect.Type) (string, interface{}) {
	for moduleName, injector := range wj.injectors {
		method := reflect.ValueOf(injector)
		methodType := method.Type()
		if methodType.NumOut() > 0 &&
			methodType.Out(0).Name() == dependencyType.Name() &&
			methodType.Out(0).PkgPath() == dependencyType.PkgPath() {
			return moduleName, injector
		}
	}

	return "", nil
}

func (wj *WireJacket) getParamTypeList(methodType reflect.Type) []reflect.Type {
	typeList := []reflect.Type{}
	for i := 0; i < methodType.NumIn(); i++ {
		dependency := methodType.In(i)
		typeList = append(typeList, dependency)
	}
	return typeList
}

func (wj *WireJacket) findModule(dependencyType reflect.Type) Module {
	for _, module := range wj.modules {
		moduleValue := reflect.ValueOf(module)
		if moduleValue.CanConvert(dependencyType) {
			return module
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
			method := reflect.ValueOf(wj.injectors[moduleName])
			methodType := method.Type()

			dependencies, satisfied := wj.getDependencies(methodType)
			if satisfied {
				returnVal := method.Call(dependencies)
				module, err := wj.checkInjectionResult(returnVal)
				if err != nil {
					return fmt.Errorf("[%s] %s", moduleName, err)
				}
				wj.modules[moduleName] = module
				activatedList = append(activatedList, moduleName)
			}
		}
		for _, activated := range activatedList {
			NotActivatedList = RemoveElement(NotActivatedList, activated)
		}
		tryCount++
	}

	return nil
}

// SetInjectors
// Wire has two basic concepts: providers and injectors.
// WireJacket's injectors stores module_name with injector as a key-value.
// The implment_name can be found in the config file.
// By default, WireJacket trys to find module_name in
// {process_name}_activating_modules of {process_name}.conf file.
//
// Example of ossicones process :
//
// # in ossicones.conf
//
// ossicones_activating_modules=ossiconesblockchain viperconfig
//
// # in wire.go
//
// func InjectViperConfig() (config.Config, error) { ... }
// func InjectOssiconesBlockchain(config config.Config) (blockchain.Blockchain, error) { ... }
//
// # injectors can be like this.
//
//	var injectors = map[string]interface{}{
// 		"viperconfig": 			InjectViperConfig,
// 		"ossiconesblockchain":	InjectOssiconesBlockchain,
// 	}
//
func (wj *WireJacket) SetInjectors(injectors map[string]interface{}) {
	// TODO
}

// SetEagerInjectors
func (wj *WireJacket) SetEagerInjectors(injectors map[string]interface{}) {
	// TODO
}

// AddInjector
func (wj *WireJacket) AddInjector(moduleName string, injector interface{}) {
	// TODO
}

// AddEagerInjector
func (wj *WireJacket) AddEagerInjector(moduleName string, injector interface{}) {
	// TODO
}

// DoWire
func (wj *WireJacket) DoWire() {
	// TODO
}

// GetConfig
func (wj *WireJacket) GetConfig() Config {
	return wj.config
}

// GetModule
func (wj *WireJacket) GetModule(moduleName string) interface{} {
	return wj.modules[moduleName]
}

// GetModuleByType
func (wj *WireJacket) GetModuleByType(interfaceType interface{}) interface{} {
	moduleType := reflect.TypeOf(interfaceType).Elem()
	for _, module := range wj.modules {
		moduleValue := reflect.ValueOf(module)
		if moduleValue.CanConvert(moduleType) {
			return module
		}
	}

	return nil
}

func readActivatingModules(config Config) []string {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	activatingModuleNames := config.GetStringSlice(
		"ossicones_activating_modules",
		[]string{},
		// defaultActivatingModules[:], // array to slice
	)

	return activatingModuleNames
}

// Close closes all the modules gracefully
func (wj *WireJacket) Close() error {
	// TODO
	return nil
}
