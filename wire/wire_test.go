package wire

import (
	"testing"

	"github.com/bang9211/wire-jacket/internal/config"
	. "github.com/stretchr/testify/assert"
)

func TestInjectOssiconesBlockchain(t *testing.T) {
	viperConfig := config.NewViperConfig("ossicones")
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")
	ossiconesBlockchain, err := InjectOssiconesBlockchain(viperConfig)
	NotNil(t, ossiconesBlockchain)
	NoError(t, err, "Failed to InjectOssiconesBlockchain()")
	Implements(t, (*Blockchain)(nil), ossiconesBlockchain, "It must implements of interface Blockchain")

	Len(t, ossiconesBlockchain.AllBlocks(), 1)
	ossiconesBlockchain.AddBlock("TEST BLOCK")
	Len(t, ossiconesBlockchain.AllBlocks(), 2)
	_, err = ossiconesBlockchain.GetBlock(1)
	NoError(t, err)
	_, err = ossiconesBlockchain.GetBlock(2)
	NoError(t, err)
	_, err = ossiconesBlockchain.GetBlock(3)
	Error(t, err)
	ossiconesBlockchain.PrintBlock()
	NoError(t, ossiconesBlockchain.Reset())
	NoError(t, ossiconesBlockchain.Close())
}

func TestInjectDefaultExplorerServer(t *testing.T) {
	viperConfig := config.NewViperConfig("ossicones")
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")
	ossiconesBlockchain, err := InjectOssiconesBlockchain(viperConfig)
	NotNil(t, ossiconesBlockchain)
	NoError(t, err, "Failed to InjectOssiconesBlockchain()")
	Implements(t, (*Blockchain)(nil), ossiconesBlockchain, "It must implements of interface Blockchain")
	defaultExplorerServer, err := InjectDefaultExplorerServer(viperConfig, ossiconesBlockchain)
	NotNil(t, defaultExplorerServer)
	NoError(t, err, "Failed to InjectDefaultExplorerServer()")
	Implements(t, (*ExplorerServer)(nil), defaultExplorerServer, "It must implements of interface ExplorerServer")

	defaultExplorerServer.Serve()
	Len(t, defaultExplorerServer.GetAllBlocks(), 1)
	NoError(t, defaultExplorerServer.Close())
}

func TestInjectDefaultRESTAPIServer(t *testing.T) {
	viperConfig := config.NewViperConfig("ossicones")
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")
	ossiconesBlockchain, err := InjectOssiconesBlockchain(viperConfig)
	NotNil(t, ossiconesBlockchain)
	NoError(t, err, "Failed to InjectOssiconesBlockchain()")
	Implements(t, (*Blockchain)(nil), ossiconesBlockchain, "It must implements of interface Blockchain")
	defaultRESTAPIServer, err := InjectDefaultRESTAPIServer(viperConfig, ossiconesBlockchain)
	NotNil(t, defaultRESTAPIServer)
	NoError(t, err, "Failed to InjectDefaultRESTAPIServer()")
	Implements(t, (*RESTAPIServer)(nil), defaultRESTAPIServer, "It must implements of interface RESTAPIServer")

	defaultRESTAPIServer.Serve()
	Equal(t, defaultRESTAPIServer.GetPaths(), []string{"/"})
	NoError(t, defaultRESTAPIServer.Close())
}
