package wirejacket

import (
	"io"
	"os"
	"testing"

	"github.com/bang9211/wire-jacket/internal/config"
	"github.com/bang9211/wire-jacket/wire"
	"github.com/stretchr/testify/assert"
)

var emptyInjectors = map[string]interface{}{}
var emptyEagerInjectors = map[string]interface{}{}

func TestWireJacketDefaultConfigCase(t *testing.T) {
	// test.conf(.envfile)
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestWireJacketSpecifiedConfigCase(t *testing.T) {
	os.Args = append(os.Args, "--config")
	os.Args = append(os.Args, "test.json")
	// test.json
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestWireJacketNoConfigCase(t *testing.T) {
	wj, err := New("no_exist_service")
	assert.NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNew(t *testing.T) {
	wj, err := New("test")
	assert.NoError(t, err, "Failed to New()")

	// no wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewNoConfigCase(t *testing.T) {
	wj, err := New("no_exist_service")
	assert.NoError(t, err, "Failed to New()")

	// no wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectors(t *testing.T) {
	wj, err := NewWithInjectors("test", emptyInjectors, emptyEagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	// no wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectorsNoConfigCase(t *testing.T) {
	wj, err := NewWithInjectors("no_exist_service", emptyInjectors, emptyEagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	// no wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetActivatingModules(t *testing.T) {
	wj, err := New("no_exist_service")
	assert.NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	// no activating modules to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.SetActivatingModules([]string{
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	// no dependency to wire
	assert.Error(t, err)

	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetInjectors(t *testing.T) {
	wj, err := New("no_exist_service")
	assert.NoError(t, err, "Failed to New()")

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)

	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	// no dependency to wire
	assert.Error(t, err)

	wj.SetInjectors(wire.Injectors)
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestSetEagerInjectors(t *testing.T) {
	wj, err := New("test")
	assert.NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.SetEagerInjectors(wire.EagerInjectors)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestAddInjector(t *testing.T) {
	wj, err := New("test")
	assert.NoError(t, err, "Failed to New()")

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	// no dependency to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestAddEagerInjector(t *testing.T) {
	wj, err := New("test")
	assert.NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestDoWire(t *testing.T) {
	wj, err := New("no_exist_service")
	assert.NoError(t, err, "Failed to New()")

	// no wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddInjector("ossiconesblockchain", wire.InjectOssiconesBlockchain)
	// no eager wire.Injectors to wire
	err = wj.DoWire()
	assert.Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", wire.InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", wire.InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetConfig(t *testing.T) {
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	config := wj.GetConfig()
	assert.Equal(t, "defaultVal", config.GetString("no_exists", "defaultVal"))
	assert.Equal(t, "Genesis OssiconesBlock", config.GetString("ossicones_blockchain_genesis_block_data", "defaultVal"))

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetModule(t *testing.T) {
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	viperconfig := wj.GetModule("viperconfig")
	assert.NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModule("ossiconesblockchain")
	assert.NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModule("defaultexplorerserver")
	assert.NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModule("defaultrestapiserver")
	assert.NotNil(t, defaultrestapiserver)
	assert.Nil(t, wj.GetModule("no_exists"))

	wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModule("viperconfig")
	assert.NotNil(t, vc)
	obc := wj.GetModule("ossiconesblockchain")
	assert.NotNil(t, obc)
	des := wj.GetModule("defaultexplorerserver")
	assert.NotNil(t, des)
	drs := wj.GetModule("defaultrestapiserver")
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModule("no_exists"))

	assert.Equal(t, viperconfig, vc)
	assert.Equal(t, ossiconesblockchain, obc)
	assert.Equal(t, defaultexplorerserver, des)
	assert.Equal(t, defaultrestapiserver, drs)

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestGetModuleByType(t *testing.T) {
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	viperconfig := wj.GetModuleByType((*config.Config)(nil))
	assert.NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModuleByType((*wire.Blockchain)(nil))
	assert.NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModuleByType((*wire.ExplorerServer)(nil))
	assert.NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModuleByType((*wire.RESTAPIServer)(nil))
	assert.NotNil(t, defaultrestapiserver)
	assert.Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	assert.Nil(t, wj.GetModuleByType(nil))

	wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModuleByType((*config.Config)(nil))
	assert.NotNil(t, vc)
	obc := wj.GetModuleByType((*wire.Blockchain)(nil))
	assert.NotNil(t, obc)
	des := wj.GetModuleByType((*wire.ExplorerServer)(nil))
	assert.NotNil(t, des)
	drs := wj.GetModuleByType((*wire.RESTAPIServer)(nil))
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	assert.Nil(t, wj.GetModuleByType(nil))

	assert.Equal(t, viperconfig, vc)
	assert.Equal(t, ossiconesblockchain, obc)
	assert.Equal(t, defaultexplorerserver, des)
	assert.Equal(t, defaultrestapiserver, drs)

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}

func TestClose(t *testing.T) {
	wj, err := NewWithInjectors("test", wire.Injectors, wire.EagerInjectors)
	assert.NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")

	err = wj.DoWire()
	assert.NoError(t, err, "Failed to DoWire()")

	viperconfig := wj.GetModule("viperconfig")
	assert.NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModule("ossiconesblockchain")
	assert.NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModule("defaultexplorerserver")
	assert.NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModule("defaultrestapiserver")
	assert.NotNil(t, defaultrestapiserver)
	assert.Nil(t, wj.GetModule("no_exists"))

	vc := wj.GetModule("viperconfig")
	assert.NotNil(t, vc)
	obc := wj.GetModule("ossiconesblockchain")
	assert.NotNil(t, obc)
	des := wj.GetModule("defaultexplorerserver")
	assert.NotNil(t, des)
	drs := wj.GetModule("defaultrestapiserver")
	assert.NotNil(t, drs)
	assert.Nil(t, wj.GetModule("no_exists"))

	err = wj.Close()
	assert.NoError(t, err, "Failed to Close()")
}
