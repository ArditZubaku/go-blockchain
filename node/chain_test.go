package node

import (
	"testing"

	"github.com/ArditZubaku/blockchain/store"
	"github.com/ArditZubaku/blockchain/types"
	"github.com/ArditZubaku/blockchain/util"
	"github.com/stretchr/testify/assert"
)

func TestAddBlock(t *testing.T) {
	chain := NewChain(store.NewMemoryBlockStore())

	for i := range 100 {
		block := util.RandomBlock()
		blockHash, err := types.HashBlock(block)
		assert.Nil(t, err)

		err = chain.AddBlock(block)
		assert.Nil(t, err)

		fetchedBlock, err := chain.GetBlockByHash(blockHash)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlock)

		fetchedBlockByHeight, err := chain.GetBlockByHeight(i)
		assert.Nil(t, err)
		assert.Equal(t, block, fetchedBlockByHeight)
	}
}

func TestChainHeight(t *testing.T) {
	chain := NewChain(store.NewMemoryBlockStore())
	for i := range 100 {
		b := util.RandomBlock()
		assert.Nil(t, chain.AddBlock(b))
		assert.Equal(t, chain.Height(), i)
	}
}
