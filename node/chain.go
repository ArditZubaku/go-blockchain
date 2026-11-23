package node

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/ArditZubaku/blockchain/proto"
	"github.com/ArditZubaku/blockchain/store"
	"github.com/ArditZubaku/blockchain/types"
)

type HeaderList struct {
	headers []*proto.Header
}

func NewHeaderList() *HeaderList {
	// TODO: Doesn't have to be 10, but why not
	return &HeaderList{headers: make([]*proto.Header, 0, 10)}
}

func (l *HeaderList) Add(h *proto.Header) {
	l.headers = append(l.headers, h)
}

func (l *HeaderList) Get(idx int) (*proto.Header, error) {
	if idx > l.Height() {
		return nil, errors.New("index too high!")
	}

	return l.headers[idx], nil
}

func (l *HeaderList) Height() int {
	return l.Len() - 1
}

// Len returns the len of the headers
// [A,B,C,D] -> len=5, height=4
func (l *HeaderList) Len() int {
	return len(l.headers)
}

type Chain struct {
	blockStore store.BlockStorer
	headers    *HeaderList
}

func NewChain(bs store.BlockStorer) *Chain {
	return &Chain{blockStore: bs, headers: NewHeaderList()}
}

func (c *Chain) Height() int {
	return c.headers.Height()
}

func (c *Chain) AddBlock(b *proto.Block) error {
	c.headers.Add(b.Header)
	return c.blockStore.Put(b)
}

func (c *Chain) GetBlockByHeight(height int) (*proto.Block, error) {
	if c.Height() < height {
		return nil, fmt.Errorf("given height %d is too high, height - %d", height, c.Height())
	}

	header, err := c.headers.Get(height)
	if err != nil {
		return nil, err
	}

	headerHash, err := types.HashHeader(header)
	if err != nil {
		return nil, err
	}

	return c.GetBlockByHash(headerHash)
}

func (c *Chain) GetBlockByHash(hash []byte) (*proto.Block, error) {
	hashHex := hex.EncodeToString(hash)
	return c.blockStore.Get(hashHex)
}
