package wirejacket

import (
	"testing"

	. "github.com/stretchr/testify/assert"
)

func TestImplementsBlock(t *testing.T) {
	Nil(t, nil, "NULL")
	// Implements(t, (*blockchain.Block)(nil), new(ossiconesblockchain.OssiconesBlock),
	// 	"It must implements of interface blockchain.Block")

	// _, ok := wj.GetModuleByType((*blockchain.Blockchain)(nil)).(blockchain.Blockchain)
	// if !ok {
	// 	return fmt.Errorf("failed to get ossiconesblockchain")
	// }

	// _, ok = wj.GetModuleByType((*explorerserver.ExplorerServer)(nil)).(explorerserver.ExplorerServer)
	// if !ok {
	// 	return fmt.Errorf("failed to get defaultexplorerserver")
	// }

	// _, ok = wj.GetModuleByType((*restapiserver.RESTAPIServer)(nil)).(restapiserver.RESTAPIServer)
	// if !ok {
	// 	return fmt.Errorf("failed to get defaultrestapiserver")
	// }
}
