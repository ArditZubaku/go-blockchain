package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePrivateKey(t *testing.T) {
	privateKey := GeneratePrivateKey()
	assert.Equal(t, len(privateKey.Bytes()), privateKeyLen)

	publicKey := privateKey.Public()
	assert.Equal(t, len(publicKey.Bytes()), publicKeyLen)
}

func TestPrivateKeySign(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.Public()
	msg := []byte("test msg")

	signature := privateKey.Sign(msg)
	assert.True(t, signature.Verify(publicKey, msg))

	// Test with invalid msg
	assert.False(t, signature.Verify(publicKey, []byte("invalid msg")))

	// Test with invalid (another) public key
	assert.False(t, signature.Verify(GeneratePrivateKey().Public(), msg))
}

func TestNewPrivateKeyFromString(t *testing.T) {
	var (
		seed       = "de7509ff4297b6ed28f839a3a86c6173e8084171617a0ffa2bcd53c7f20381c7"
		privateKey = NewPrivateKeyFromString(seed)
		addressStr = "503c28fd1508de1540a8f55fdf20b2d7dcfe5b25"
	)

	assert.Equal(t, privateKeyLen, len(privateKey.Bytes()))

	address := privateKey.Public().Address()
	assert.Equal(t, addressStr, address.String())
}

func TestPublicKeyToAddress(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.Public()
	address := publicKey.Address()

	assert.Equal(t, addressLen, len(address.Bytes()))
}
