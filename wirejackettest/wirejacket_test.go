package wirejackettest

import (
	"io"
	"os"
	"testing"

	wirejacket "github.com/bang9211/wire-jacket"
	"github.com/bang9211/wire-jacket/internal/config"
	. "github.com/stretchr/testify/assert"
)

var emptyInjectors = map[string]interface{}{}
var emptyEagerInjectors = map[string]interface{}{}

func TestWireJacketDefaultConfigCase(t *testing.T) {
	// ossicones.conf(.envfile)
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestWireJacketSpecifiedConfigCase(t *testing.T) {
	os.Args = append(os.Args, "--config")
	os.Args = append(os.Args, "ossicones.json")
	// ossicones.json
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestWireJacketNoConfigCase(t *testing.T) {
	wj, err := wirejacket.New("no_exist_service")
	NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestNew(t *testing.T) {
	wj, err := wirejacket.New("ossicones")
	NoError(t, err, "Failed to New()")

	// no injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestNewNoConfigCase(t *testing.T) {
	wj, err := wirejacket.New("no_exist_service")
	NoError(t, err, "Failed to New()")

	// no injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectors(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("ossicones", emptyInjectors, emptyEagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	// no injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestNewWithEmptyInjectorsNoConfigCase(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("no_exist_service", emptyInjectors, emptyEagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	// no injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestSetActivatingModules(t *testing.T) {
	wj, err := wirejacket.New("no_exist_service")
	NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	// no activating modules to wire
	err = wj.DoWire()
	Error(t, err)

	wj.SetActivatingModules([]string{
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	// no dependency to wire
	Error(t, err)

	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestSetInjectors(t *testing.T) {
	wj, err := wirejacket.New("no_exist_service")
	NoError(t, err, "Failed to New()")

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)

	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	// no dependency to wire
	Error(t, err)

	wj.SetInjectors(injectors)
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestSetEagerInjectors(t *testing.T) {
	wj, err := wirejacket.New("ossicones")
	NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.SetEagerInjectors(eagerInjectors)

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestAddInjector(t *testing.T) {
	wj, err := wirejacket.New("ossicones")
	NoError(t, err, "Failed to New()")

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	// no dependency to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestAddEagerInjector(t *testing.T) {
	wj, err := wirejacket.New("ossicones")
	NoError(t, err, "Failed to New()")

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestDoWire(t *testing.T) {
	wj, err := wirejacket.New("no_exist_service")
	NoError(t, err, "Failed to New()")

	// no injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddInjector("ossiconesblockchain", InjectOssiconesBlockchain)
	// no eager injectors to wire
	err = wj.DoWire()
	Error(t, err)

	wj.AddEagerInjector("defaultexplorerserver", InjectDefaultExplorerServer)
	wj.AddEagerInjector("defaultrestapiserver", InjectDefaultRESTAPIServer)
	wj.SetActivatingModules([]string{
		"ossiconesblockchain",
		"defaultexplorerserver",
		"defaultrestapiserver",
	})
	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestGetConfig(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	config := wj.GetConfig()
	Equal(t, "defaultVal", config.GetString("no_exists", "defaultVal"))
	Equal(t, "Genesis OssiconesBlock", config.GetString("ossicones_blockchain_genesis_block_data", "defaultVal"))

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestGetModule(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	viperconfig := wj.GetModule("viperconfig")
	NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModule("ossiconesblockchain")
	NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModule("defaultexplorerserver")
	NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModule("defaultrestapiserver")
	NotNil(t, defaultrestapiserver)
	Nil(t, wj.GetModule("no_exists"))

	wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModule("viperconfig")
	NotNil(t, vc)
	obc := wj.GetModule("ossiconesblockchain")
	NotNil(t, obc)
	des := wj.GetModule("defaultexplorerserver")
	NotNil(t, des)
	drs := wj.GetModule("defaultrestapiserver")
	NotNil(t, drs)
	Nil(t, wj.GetModule("no_exists"))

	Equal(t, viperconfig, vc)
	Equal(t, ossiconesblockchain, obc)
	Equal(t, defaultexplorerserver, des)
	Equal(t, defaultrestapiserver, drs)

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestGetModuleByType(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	viperconfig := wj.GetModuleByType((*config.Config)(nil))
	NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModuleByType((*Blockchain)(nil))
	NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModuleByType((*ExplorerServer)(nil))
	NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModuleByType((*RESTAPIServer)(nil))
	NotNil(t, defaultrestapiserver)
	Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	Nil(t, wj.GetModuleByType(nil))

	wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	vc := wj.GetModuleByType((*config.Config)(nil))
	NotNil(t, vc)
	obc := wj.GetModuleByType((*Blockchain)(nil))
	NotNil(t, obc)
	des := wj.GetModuleByType((*ExplorerServer)(nil))
	NotNil(t, des)
	drs := wj.GetModuleByType((*RESTAPIServer)(nil))
	NotNil(t, drs)
	Nil(t, wj.GetModuleByType((*io.Writer)(nil)))
	Nil(t, wj.GetModuleByType(nil))

	Equal(t, viperconfig, vc)
	Equal(t, ossiconesblockchain, obc)
	Equal(t, defaultexplorerserver, des)
	Equal(t, defaultrestapiserver, drs)

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}

func TestClose(t *testing.T) {
	wj, err := wirejacket.NewWithInjectors("ossicones", injectors, eagerInjectors)
	NoError(t, err, "Failed to NewWithInjectors()")

	err = wj.Close()
	NoError(t, err, "Failed to Close()")

	err = wj.DoWire()
	NoError(t, err, "Failed to DoWire()")

	viperconfig := wj.GetModule("viperconfig")
	NotNil(t, viperconfig)
	ossiconesblockchain := wj.GetModule("ossiconesblockchain")
	NotNil(t, ossiconesblockchain)
	defaultexplorerserver := wj.GetModule("defaultexplorerserver")
	NotNil(t, defaultexplorerserver)
	defaultrestapiserver := wj.GetModule("defaultrestapiserver")
	NotNil(t, defaultrestapiserver)
	Nil(t, wj.GetModule("no_exists"))

	vc := wj.GetModule("viperconfig")
	NotNil(t, vc)
	obc := wj.GetModule("ossiconesblockchain")
	NotNil(t, obc)
	des := wj.GetModule("defaultexplorerserver")
	NotNil(t, des)
	drs := wj.GetModule("defaultrestapiserver")
	NotNil(t, drs)
	Nil(t, wj.GetModule("no_exists"))

	err = wj.Close()
	NoError(t, err, "Failed to Close()")
}
