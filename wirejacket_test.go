package wirejacket

import (
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/bang9211/wire-jacket/internal/mockup"
	"github.com/stretchr/testify/assert"
)

var emptyInjectors = map[string]interface{}{}
var emptyEagerInjectors = map[string]interface{}{}

func TestNewDefaultConfigCase(t *testing.T) {
	// app.conf(.envfile)
	wj := New().
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewWithServiceNameDefaultConfigCase(t *testing.T) {
	// app.conf(.envfile)
	wj := NewWithServiceName("test example").
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestWireJacketSpecifiedConfigCase(t *testing.T) {
	backupArgs := os.Args

	os.Args = append(os.Args, "--config")
	os.Args = append(os.Args, "test.json")
	// test.json
	wj := NewWithServiceName("test example").
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")

	os.Args = backupArgs
}

func TestWireJacketNoConfigCase(t *testing.T) {
	backupArgs := os.Args

	os.Args = append(os.Args, "--config")
	os.Args = append(os.Args, "no_exist.json")
	wj := NewWithServiceName("no_exist_service")

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)
	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")

	os.Args = backupArgs
}

func TestNew(t *testing.T) {
	wj := New()

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectors(t *testing.T) {
	wj := NewWithServiceName("test example").
		SetInjectors(emptyInjectors).
		SetEagerInjectors(emptyEagerInjectors)

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectorsNoConfigCase(t *testing.T) {
	wj := NewWithServiceName("no_exist_service").
		SetInjectors(emptyInjectors).
		SetEagerInjectors(emptyEagerInjectors)

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)
	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetActivatingModules(t *testing.T) {
	wj := NewWithServiceName("no_exist_service")

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)
	// no activating modules to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.SetActivatingModules([]string{
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	// no dependency to wire
	assert.Error(t, err)

	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetInjectors(t *testing.T) {
	wj := NewWithServiceName("no_exist_service")

	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)

	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err := wj.DoWire()
	// no dependency to wire
	assert.Error(t, err)

	wj.SetInjectors(mockup.Injectors)
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetEagerInjectors(t *testing.T) {
	wj := New()

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.SetEagerInjectors(mockup.EagerInjectors)

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestAddInjector(t *testing.T) {
	wj := New()

	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)
	// no dependency to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestAddEagerInjector(t *testing.T) {
	wj := New()

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestDoWire(t *testing.T) {
	wj := NewWithServiceName("no_exist_service")

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)

	wj.SetActivatingModules([]string{})
	err = wj.DoWire()
	assert.Error(t, err)

	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestDoWireWithoutEagerInjectors(t *testing.T) {
	wj := NewWithServiceName("no_exist_service")

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddInjector("mockup_restapiserver", mockup.InjectMockupRESTAPIServer)
	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetConfig(t *testing.T) {
	wj := New().
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	config := GetConfig()
	assert.Equal(t, "defaultVal", config.GetString("no_exists", "defaultVal"))
	assert.Equal(t, "Genesis OssiconesBlock", config.GetString("ossicones_genesis_block_data", "defaultVal"))

	err := wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetModule(t *testing.T) {
	wj := New().
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	viperconfig := wj.GetModule("viperconfig")
	assert.NotNil(t, viperconfig)
	mockup_blockchain := wj.GetModule("mockup_blockchain")
	assert.NotNil(t, mockup_blockchain)
	mockup_explorerserver := wj.GetModule("mockup_explorerserver")
	assert.NotNil(t, mockup_explorerserver)
	mockup_restapiserver := wj.GetModule("mockup_restapiserver")
	assert.NotNil(t, mockup_restapiserver)
	assert.Nil(t, wj.GetModule("no_exists"))

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModule("viperconfig")
	assert.NotNil(t, vc)
	obc := wj.GetModule("mockup_blockchain")
	assert.NotNil(t, obc)
	des := wj.GetModule("mockup_explorerserver")
	assert.NotNil(t, des)
	drs := wj.GetModule("mockup_restapiserver")
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModule("no_exists"))

	assert.Equal(t, viperconfig, vc)
	assert.Equal(t, mockup_blockchain, obc)
	assert.Equal(t, mockup_explorerserver, des)
	assert.Equal(t, mockup_restapiserver, drs)

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetModuleByType(t *testing.T) {
	wj := New().
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	viperconfig := wj.GetModuleByType((*config.Config)(nil))
	assert.NotNil(t, viperconfig)
	mockup_blockchain := wj.GetModuleByType((*mockup.Blockchain)(nil))
	assert.NotNil(t, mockup_blockchain)
	mockup_explorerserver := wj.GetModuleByType((*mockup.ExplorerServer)(nil))
	assert.NotNil(t, mockup_explorerserver)
	mockup_restapiserver := wj.GetModuleByType((*mockup.RESTAPIServer)(nil))
	assert.NotNil(t, mockup_restapiserver)
	assert.Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	assert.Nil(t, wj.GetModuleByType(nil))

	err := wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModuleByType((*config.Config)(nil))
	assert.NotNil(t, vc)
	obc := wj.GetModuleByType((*mockup.Blockchain)(nil))
	assert.NotNil(t, obc)
	des := wj.GetModuleByType((*mockup.ExplorerServer)(nil))
	assert.NotNil(t, des)
	drs := wj.GetModuleByType((*mockup.RESTAPIServer)(nil))
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	assert.Nil(t, wj.GetModuleByType(nil))

	assert.Equal(t, viperconfig, vc)
	assert.Equal(t, mockup_blockchain, obc)
	assert.Equal(t, mockup_explorerserver, des)
	assert.Equal(t, mockup_restapiserver, drs)

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestDoWireFailedInjectorCall(t *testing.T) {
	wj := New()

	// no mockup.Injectors to wire
	err := wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("mockup_database", mockup.InjectMockupDB)
	wj.AddInjector("mockup_blockchain", mockup.InjectMockupBlockchain)
	wj.AddEagerInjector("mockup_explorerserver", mockup.InjectMockupExplorerServer)
	wj.AddEagerInjector("mockup_restapiserver", mockup.InjectMockupInvalidReturnTest)

	wj.SetActivatingModules([]string{
		"mockup_database",
		"mockup_blockchain",
		"mockup_explorerserver",
		"mockup_restapiserver",
	})
	err = wj.DoWire()
	assert.Error(t, err)

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestCheckInjectionResult(t *testing.T) {
	wj := New()
	values := []reflect.Value{reflect.Value{}}
	_, err := wj.checkInjectionResult(values)
	assert.Error(t, err)

	values = []reflect.Value{reflect.ValueOf(true)}
	_, err = wj.checkInjectionResult(values)
	assert.Error(t, err)

	values = []reflect.Value{reflect.Value{}, reflect.ValueOf(true)}
	_, err = wj.checkInjectionResult(values)
	assert.Error(t, err)

	// get dependencies
	injectorFunc := reflect.ValueOf(mockup.InjectMockupDB)
	dependencies, err := wj.getDependencies("mockup_database", injectorFunc.Type())

	// call injector
	returnVal := injectorFunc.Call(dependencies)

	values = []reflect.Value{reflect.Value{}, reflect.Value{}}
	_, err = wj.checkInjectionResult(values)
	assert.Error(t, err)

	values = []reflect.Value{reflect.Value{}, returnVal[1]}
	_, err = wj.checkInjectionResult(values)
	assert.Error(t, err)

	values = []reflect.Value{reflect.ValueOf(true), returnVal[1]}
	_, err = wj.checkInjectionResult(values)
	assert.Error(t, err)
}

func TestClose(t *testing.T) {
	wj := New().
		SetInjectors(mockup.Injectors).
		SetEagerInjectors(mockup.EagerInjectors)

	err := wj.Close()
	assert.NoError(t, err, "Failed to Close()")

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	viperconfig := wj.GetModule("viperconfig")
	assert.NotNil(t, viperconfig)
	mockup_blockchain := wj.GetModule("mockup_blockchain")
	assert.NotNil(t, mockup_blockchain)
	mockup_explorerserver := wj.GetModule("mockup_explorerserver")
	assert.NotNil(t, mockup_explorerserver)
	mockup_restapiserver := wj.GetModule("mockup_restapiserver")
	assert.NotNil(t, mockup_restapiserver)
	assert.Nil(t, wj.GetModule("no_exists"))

	vc := wj.GetModule("viperconfig")
	assert.NotNil(t, vc)
	obc := wj.GetModule("mockup_blockchain")
	assert.NotNil(t, obc)
	des := wj.GetModule("mockup_explorerserver")
	assert.NotNil(t, des)
	drs := wj.GetModule("mockup_restapiserver")
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModule("no_exists"))

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestCloseErro(t *testing.T) {
	wj := New()
	wj.AddEagerInjector("test", mockup.InjectMockupInvalidImplTest)
	assert.NoError(t, wj.Close())

	wj.SetActivatingModules([]string{"test"})
	wj.GetModule("test")
	assert.NoError(t, wj.Close())
}
