package types

import (
	"crypto/sha256"

	"github.com/ArditZubaku/blockchain/crypto"
	"github.com/ArditZubaku/blockchain/proto"
	pb "google.golang.org/protobuf/proto"
)

func SignTransaction(pk crypto.PrivateKey, tx *proto.Transaction) crypto.Signature {
	return pk.Sign(HashTransaction(tx))
}

func HashTransaction(tx *proto.Transaction) []byte {
	b, err := pb.Marshal(tx)
	if err != nil {
		panic(err)
	}

	hash := sha256.Sum256(b)

	return hash[:]
}

func VerifyTransaction(tx *proto.Transaction) bool {
	for _, input := range tx.Inputs {
		var (
			signature = crypto.NewSignatureFromBytes(input.Signature)
			publicKey = crypto.NewPublicKeyFromBytes(input.PublicKey)
		)
		input.Signature = nil
		if !signature.Verify(publicKey, HashTransaction(tx)) {
			return false
		}
	}

	return true
}
