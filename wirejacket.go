package wirejacket

import (
	"fmt"
	"log"
	"reflect"

	"github.com/bang9211/ossicones/utils"
)

type modulable interface {
	Close() error
}

type WireJacket struct {
	injection_list map[string]interface{}
	instance_list  map[string]modulable
	config         Config
}

func New(injection_list map[string]interface{}) (*WireJacket, error) {
	wj := &WireJacket{
		injection_list: injection_list,
		instance_list:  map[string]modulable{},
		config:         NewViperConfig(),
	}

	activatingModules := wj.readActivatingModules(wj.config)
	NotActivatedList := make([]string, len(activatingModules))
	copy(NotActivatedList, activatingModules)

	activatedList := []string{}
	tryCount := 0
	for len(NotActivatedList) > 0 && tryCount < len(NotActivatedList)*len(NotActivatedList) {
		for _, moduleName := range NotActivatedList {
			method := reflect.ValueOf(injection_list[moduleName])
			methodType := method.Type()

			dependencies, satisfied := wj.getNecessaryDependencies(methodType)
			if satisfied {
				returnVal := method.Call(dependencies)
				modulableModule, err := wj.checkInjectionResult(returnVal)
				if err != nil {
					return nil, fmt.Errorf("[%s] %s", moduleName, err)
				}
				wj.instance_list[moduleName] = modulableModule
				activatedList = append(activatedList, moduleName)
			}
		}
		for _, activated := range activatedList {
			NotActivatedList = utils.RemoveElement(NotActivatedList, activated)
		}
		tryCount++
	}

	return wj, nil
}

func (wj *WireJacket) GetConfig() Config {
	return wj.config
}

func (wj *WireJacket) GetInstance(moduleName string) interface{} {
	return wj.instance_list[moduleName]
}

func (wj *WireJacket) GetInstanceByType(interfaceType interface{}) interface{} {
	return wj.instance_list
}

func (wj *WireJacket) readActivatingModules(config Config) []string {
	err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	activatingModules := config.GetStringSlice(
		"ossicones_activating_modules",
		[]string{},
		// defaultActivatingModules[:], // array to slice
	)

	return activatingModules
}

func (wj *WireJacket) getNecessaryDependencies(methodType reflect.Type) ([]reflect.Value, bool) {
	dependencies := []reflect.Value{}
	for i := 0; i < methodType.NumIn(); i++ {
		dependency := methodType.In(i)
		find := false
		for _, instance := range wj.instance_list {
			instanceValue := reflect.ValueOf(instance)
			if instanceValue.CanConvert(dependency) {
				dependencies = append(dependencies, instanceValue)
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

func (wj *WireJacket) checkInjectionResult(returnVal []reflect.Value) (modulable, error) {

	if len(returnVal) != 1 && len(returnVal) != 2 {
		return nil, fmt.Errorf(
			"invalid inject function format len(return) : %d", len(returnVal))
	}
	var modulableModule modulable
	var ok bool
	if len(returnVal) == 1 { // return (instance)
		if !returnVal[0].CanInterface() {
			return nil, fmt.Errorf(
				"returnVal(%s) can't be interface",
				returnVal[0],
			)
		}
		modulableModule, ok = returnVal[0].Interface().(modulable)
		if !ok {
			return nil, fmt.Errorf(
				"failed to cast returnVal(%s) to modulable", returnVal[0])
		}
	} else { // return (instance, error)
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
		modulableModule, ok = returnVal[0].Interface().(modulable)
		if !ok {
			return nil, fmt.Errorf(
				"failed to cast returnVal(%s) to modulable", returnVal[0])
		}
	}
	return modulableModule, nil
}
