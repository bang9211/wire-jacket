package wire

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"

	"github.com/bang9211/wire-jacket/internal/config"
)

const defaultGenesisBlockData = "Genesis OssiconesBlock"

var obc *OssiconesBlockchain
var obcOnce sync.Once

type Block interface {
	// CalculateHash calculates hash using sha256.
	CalculateHash()
	// GetData gets data of the block.
	GetData() string
}

type Blockchain interface {
	// AddBlock adds data to blockchain.
	AddBlock(data string)
	// AllBlocks gets all the blocks of this blockchain.
	AllBlocks() []interface{}
	// PrintBlock just prints all the blocks.
	PrintBlock()
	// GetBlock get block at the height of this blockchain.
	GetBlock(hegiht int) (Block, error)
	// Reset resets blockchain data.
	Reset() error
	// Close closes blockchain.
	Close() error
}

var ErrorNotFound = errors.New("block not found")
var ErrorUnknown = errors.New("unknown")

// OssiconesBlock for OssiconesBlockChain.
type OssiconesBlock struct {
	Data     string `json:"data"`
	Hash     string `json:"hash"`
	PrevHash string `json:"prevhash,omitempty"`
	Height   int    `json:"height"`
}

func (b *OssiconesBlock) CalculateHash() {
	Hash := sha256.Sum256([]byte(b.Data + b.PrevHash))
	b.Hash = fmt.Sprintf("%x", Hash)
}

func (b *OssiconesBlock) GetData() string {
	return b.Data
}

type OssiconesBlockchain struct {
	config config.Config
	blocks []*OssiconesBlock
}

// GetOrCreate returns the existing singletone object of OssiconesBlockchain if present.
// Otherwise, it creates and returns the object.
func GetOrCreateOssiconesBlockchain(config config.Config) Blockchain {
	if obc == nil {
		obcOnce.Do(func() {
			obc = &OssiconesBlockchain{config: config}
			data := config.GetString(
				"OSSICONES_BLOCKCHAIN_GENESIS_BLOCK_DATA",
				defaultGenesisBlockData)
			obc.AddBlock(data)
		})
	}
	return obc
}

func (o *OssiconesBlockchain) createBlock(Data string) *OssiconesBlock {
	newBlock := OssiconesBlock{
		Data:     Data,
		PrevHash: o.getLastBlockHash(),
		Height:   len(o.blocks) + 1,
	}
	newBlock.CalculateHash()
	return &newBlock
}

func (o *OssiconesBlockchain) getLastBlockHash() string {
	if len(o.blocks) > 0 {
		return o.blocks[len(o.blocks)-1].Hash
	}
	return ""
}

func (o *OssiconesBlockchain) AddBlock(Data string) {
	newBlock := o.createBlock(Data)
	newBlock.CalculateHash()
	o.blocks = append(o.blocks, newBlock)
}

func (o *OssiconesBlockchain) AllBlocks() []interface{} {
	blocks := make([]interface{}, len(o.blocks))
	for i, block := range o.blocks {
		blocks[i] = block
	}
	return blocks
}

func (o *OssiconesBlockchain) GetBlock(height int) (Block, error) {
	if height > len(o.blocks) {
		return nil, ErrorNotFound
	}
	return o.blocks[height-1], nil
}

func (o *OssiconesBlockchain) PrintBlock() {
	for i, OssiconesBlock := range o.blocks {
		fmt.Println(i, ":", *OssiconesBlock)
	}
}

func (o *OssiconesBlockchain) Reset() error {
	o.blocks = []*OssiconesBlock{}
	data := o.config.GetString(
		"OSSICONES_BLOCKCHAIN_GENESIS_BLOCK_DATA",
		defaultGenesisBlockData)
	obc.AddBlock(data)
	return nil
}

func (o *OssiconesBlockchain) Close() error {
	return nil
}
