package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"io"
)

const (
	privateKeyLen = 64
	signatureLen  = 64
	publicKeyLen  = 32
	seedLen       = 32
	addressLen    = 20
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func NewPrivateKeyFromSeed(seed []byte) PrivateKey {
	if len(seed) != seedLen {
		panic("")
	}
	return PrivateKey{key: ed25519.NewKeyFromSeed(seed)}
}

func NewPrivateKeyFromString(s string) PrivateKey {
	seed, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}

	return NewPrivateKeyFromSeed(seed)
}

func GeneratePrivateKey() PrivateKey {
	seed := make([]byte, seedLen)
	_, err := io.ReadFull(rand.Reader, seed)
	if err != nil {
		panic("invalid seed length, must be 32")
	}

	return NewPrivateKeyFromSeed(seed)
}

func (p PrivateKey) Bytes() []byte {
	return p.key
}

func (p PrivateKey) Sign(msg []byte) Signature {
	return Signature{value: ed25519.Sign(p.key, msg)}
}

func (p PrivateKey) Public() PublicKey {
	buf := make([]byte, publicKeyLen)

	copy(buf, p.key[32:])

	return PublicKey{key: buf}
}

type PublicKey struct {
	key ed25519.PublicKey
}

func NewPublicKeyFromBytes(b []byte) PublicKey {
	if len(b) != publicKeyLen {
		panic("invalid public key length")
	}
	return PublicKey{key: ed25519.PublicKey(b)}
}

func (p PublicKey) Bytes() []byte {
	return p.key
}

func (p PublicKey) Address() Address {
	return Address{value: p.key[len(p.key)-addressLen:]}
}

type Signature struct {
	value []byte
}

func NewSignatureFromBytes(b []byte) Signature {
	if len(b) != signatureLen {
		panic("length of the bytes is not equal to signatureLen")
	}

	return Signature{value: b}
}

func (s Signature) Bytes() []byte {
	return s.value
}

func (s Signature) Verify(publicKey PublicKey, msg []byte) bool {
	return ed25519.Verify(publicKey.key, msg, s.value)
}

type Address struct {
	value []byte
}

func (a Address) String() string {
	return hex.EncodeToString(a.value)
}

func (a Address) Bytes() []byte {
	return a.value
}
