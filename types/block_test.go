package types

import (
	"testing"

	"github.com/ArditZubaku/blockchain/crypto"
	"github.com/ArditZubaku/blockchain/util"
	"github.com/stretchr/testify/assert"
)

func TestHashBlock(t *testing.T) {
	block := util.RandomBlock()
	hash := HashBlock(block)
	assert.Equal(t, 32, len(hash))
}

func TestSignBlock(t *testing.T) {
	var (
		block      = util.RandomBlock()
		privateKey = crypto.GeneratePrivateKey()
		publicKey  = privateKey.Public()
	)

	signature := SignBlock(privateKey, block)
	assert.Equal(t, 64, len(signature.Bytes()))
	assert.True(t, signature.Verify(publicKey, HashBlock(block)))
}
