package wirejackettest

import (
	"os"
	"testing"

	wirejacket "github.com/bang9211/wire-jacket"
	. "github.com/stretchr/testify/assert"
)

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

}

func TestNewWithInjectors(t *testing.T) {

}

func TestSetActivatingModules(t *testing.T) {

}

func TestSetInjectors(t *testing.T) {

}

func TestSetEagerInjectors(t *testing.T) {

}

func TestAddInjector(t *testing.T) {

}

func TestAddEagerInjector(t *testing.T) {

}

func TestDoWire(t *testing.T) {

}

func TestGetConfig(t *testing.T) {

}

func TestGetModule(t *testing.T) {

}

func TestGetModuleByType(t *testing.T) {

}

func TestClose(t *testing.T) {

}
