package store

import (
	"encoding/hex"
	"fmt"
	"sync"

	"github.com/ArditZubaku/blockchain/proto"
	"github.com/ArditZubaku/blockchain/types"
)

type BlockStorer interface {
	Put(*proto.Block) error
	Get(hash string) (*proto.Block, error)
}

type memoryBlockStore struct {
	lock   sync.RWMutex
	blocks map[string]*proto.Block
}

func NewMemoryBlockStore() BlockStorer {
	return &memoryBlockStore{
		blocks: make(map[string]*proto.Block),
	}
}

// Put implements [BlockStorer].
func (m *memoryBlockStore) Put(b *proto.Block) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	hashB, err := types.HashBlock(b)
	if err != nil {
		return err
	}

	hash := hex.EncodeToString(hashB)
	m.blocks[hash] = b
	return nil
}

// Get implements [BlockStorer].
func (m *memoryBlockStore) Get(hash string) (*proto.Block, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	block, ok := m.blocks[hash]
	if !ok {
		return nil, fmt.Errorf("block with hash %s does not exist", hash)
	}

	return block, nil
}
