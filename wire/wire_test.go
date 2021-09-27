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
}
