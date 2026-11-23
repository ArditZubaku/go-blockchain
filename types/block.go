package types

import (
	"crypto/sha256"

	"github.com/ArditZubaku/blockchain/crypto"
	"github.com/ArditZubaku/blockchain/proto"
	pb "google.golang.org/protobuf/proto"
)

// HashBlock returns a SHA256 of the header.
func HashBlock(b *proto.Block) ([]byte, error) {
	return HashHeader(b.Header)
}

func HashHeader(h *proto.Header) ([]byte, error) {
	b, err := pb.Marshal(h)
	if err != nil {
		return nil, err
	}

	hash := sha256.Sum256(b)

	return hash[:], nil
}

func SignBlock(pk crypto.PrivateKey, b *proto.Block) (crypto.Signature, error) {
	hash, err := HashBlock(b)
	if err != nil {
		return crypto.Signature{}, err
	}
	return pk.Sign(hash), nil
}
