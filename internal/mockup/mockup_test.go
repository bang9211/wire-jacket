package mockup

import (
	"testing"

	"github.com/bang9211/wire-jacket/internal/config"
	. "github.com/stretchr/testify/assert"
)

func TestInjectMockupDB(t *testing.T) {
	viperConfig := config.GetOrCreate()
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")

	mysql, err := InjectMockupDB(viperConfig)
	NotNil(t, mysql)
	NoError(t, err)
	Implements(t, (*Database)(nil), mysql, "It must implements of interface Database")

	NoError(t, mysql.Connect())
	NoError(t, mysql.Close())
}

func TestInjectMockupBlockchain(t *testing.T) {
	viperConfig := config.GetOrCreate()
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")

	mockupDB, err := InjectMockupDB(viperConfig)
	NotNil(t, mockupDB)
	NoError(t, err, "Failed to InjectMockupDB()")
	Implements(t, (*Database)(nil), mockupDB, "It must implements of interface Database")

	mockbupBlockchain, err := InjectMockupBlockchain(mockupDB)
	NotNil(t, mockbupBlockchain)
	NoError(t, err, "Failed to InjectMockupBlockchain()")
	Implements(t, (*Blockchain)(nil), mockbupBlockchain, "It must implements of interface Blockchain")

	NoError(t, mockbupBlockchain.Init(), "Failed to Init()")
	NoError(t, mockbupBlockchain.AddBlock("test data"), "Failed to AddBlock()")
	blocks := mockbupBlockchain.GetBlocks()
	Len(t, blocks, 2)
	Equal(t, blocks[0].GetData(), genesisBlockData)
	Equal(t, blocks[1].GetData(), "test data")
	NoError(t, mockbupBlockchain.Close())
}

func TestInjectMockupExplorerServer(t *testing.T) {
	viperConfig := config.GetOrCreate()
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")

	mockupDB, err := InjectMockupDB(viperConfig)
	NotNil(t, mockupDB)
	NoError(t, err, "Failed to InjectMockupDB()")
	Implements(t, (*Database)(nil), mockupDB, "It must implements of interface Database")

	mockbupBlockchain, err := InjectMockupBlockchain(mockupDB)
	NotNil(t, mockbupBlockchain)
	NoError(t, err, "Failed to InjectMockupBlockchain()")
	Implements(t, (*Blockchain)(nil), mockbupBlockchain, "It must implements of interface Blockchain")
	NoError(t, mockbupBlockchain.Init(), "Failed to Init()")

	mockupExplorerServer, err := InjectMockupExplorerServer(viperConfig, mockbupBlockchain)
	NotNil(t, mockupExplorerServer)
	NoError(t, err, "Failed to InjectMockupExplorerServer()")
	Implements(t, (*ExplorerServer)(nil), mockupExplorerServer, "It must implements of interface ExplorerServer")

	NoError(t, mockupExplorerServer.Serve())
	Len(t, mockupExplorerServer.GetAllBlockData(), 1)
	NoError(t, mockupExplorerServer.Close())
}

func TestInjectMockupRESTAPIServer(t *testing.T) {
	viperConfig := config.GetOrCreate()
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")

	mockupDB, err := InjectMockupDB(viperConfig)
	NotNil(t, mockupDB)
	NoError(t, err, "Failed to InjectMockupDB()")
	Implements(t, (*Database)(nil), mockupDB, "It must implements of interface Database")

	mockbupBlockchain, err := InjectMockupBlockchain(mockupDB)
	NotNil(t, mockbupBlockchain)
	NoError(t, err, "Failed to InjectMockupBlockchain()")
	Implements(t, (*Blockchain)(nil), mockbupBlockchain, "It must implements of interface Blockchain")
	NoError(t, mockbupBlockchain.Init(), "Failed to Init()")

	mockupRESTAPIServer, err := InjectMockupRESTAPIServer(viperConfig, mockbupBlockchain)
	NotNil(t, mockupRESTAPIServer)
	NoError(t, err, "Failed to InjectMockupRESTAPIServer()")
	Implements(t, (*RESTAPIServer)(nil), mockupRESTAPIServer, "It must implements of interface RESTAPIServer")

	NoError(t, mockupRESTAPIServer.Serve())
	Equal(t, mockupRESTAPIServer.GetPaths(), []string{"/"})
	NoError(t, mockupRESTAPIServer.Close())
}

func TestInjectMockupInvalidReturnTest(t *testing.T) {
	viperConfig := config.GetOrCreate()
	Implements(t, (*config.Config)(nil), viperConfig, "It must implements of interface config.Config")

	mockupDB, err := InjectMockupDB(viperConfig)
	NotNil(t, mockupDB)
	NoError(t, err, "Failed to InjectMockupDB()")
	Implements(t, (*Database)(nil), mockupDB, "It must implements of interface Database")

	mockbupBlockchain, err := InjectMockupBlockchain(mockupDB)
	NotNil(t, mockbupBlockchain)
	NoError(t, err, "Failed to InjectMockupBlockchain()")
	Implements(t, (*Blockchain)(nil), mockbupBlockchain, "It must implements of interface Blockchain")
	NoError(t, mockbupBlockchain.Init(), "Failed to Init()")

	testImpl, f, err := InjectMockupInvalidReturnTest(viperConfig, mockbupBlockchain)
	NotNil(t, testImpl)
	NotNil(t, f)
	f()
	NoError(t, err, "Failed to InjectMockupInvalidImplTest()")
}

func TestInjectMockupInvalidImplTestj(t *testing.T) {
	testImpl, err := InjectMockupInvalidImplTest()
	NotNil(t, testImpl)
	NoError(t, err, "Failed to InjectMockupInvalidImplTest()")
	NoError(t, testImpl.Test())
	Error(t, testImpl.Close())
}
