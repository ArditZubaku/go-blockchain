package types

import (
	"testing"

	"github.com/ArditZubaku/blockchain/crypto"
	"github.com/ArditZubaku/blockchain/proto"
	"github.com/ArditZubaku/blockchain/util"
	"github.com/stretchr/testify/assert"
)

func TestNewTransaction(t *testing.T) {
	// UTXO - Unspent Transaction Output
	// my balance is 100 coins
	// I want to send 5 coins to "test_address"
	// 2 outputs
	// 1 to the guy we wanna send to
	// 1 to ourselves

	fromPrivateKey := crypto.GeneratePrivateKey()
	fromAddress := fromPrivateKey.Public().Address().Bytes()

	toPrivateKey := crypto.GeneratePrivateKey()
	toAddress := toPrivateKey.Public().Address().Bytes()

	// input is basically info about ourselves
	input := proto.TxInput{
		PreviousTxHash: util.RandomHash(),
		PreviousOutIdx: 0,
		PublicKey:      fromPrivateKey.Public().Bytes(),
	}

	// output is specifying the destination
	output1 := proto.TxOutput{
		Amount:  5,
		Address: toAddress,
	}

	// we need to "spend" the rest as well
	output2 := proto.TxOutput{
		Amount:  5,
		Address: fromAddress, // we send it back to ourselves
	}

	tx := proto.Transaction{
		Version: 1,
		Inputs:  []*proto.TxInput{&input},
		Outputs: []*proto.TxOutput{&output1, &output2},
	}

	signature := SignTransaction(fromPrivateKey, &tx)
	input.Signature = signature.Bytes()

	assert.True(t, VerifyTransaction(&tx))
}
