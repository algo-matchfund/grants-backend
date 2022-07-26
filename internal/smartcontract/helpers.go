package smartcontract

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"fmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/encoding/msgpack"
	"github.com/algorand/go-algorand-sdk/mnemonic"
	"github.com/algorand/go-algorand-sdk/types"
)

type account struct {
	Address    types.Address
	PrivateKey ed25519.PrivateKey
}

func getAccount(m string) (*account, error) {
	privateKey, err := mnemonic.ToPrivateKey(m)

	if err != nil {
		return nil, err
	}
	publicKey := privateKey.Public()
	var a types.Address
	cpk := publicKey.(ed25519.PublicKey)
	copy(a[:], cpk[:])
	address, err := types.DecodeAddress(a.String())
	if err != nil {
		return nil, err
	}

	return &account{
		Address:    address,
		PrivateKey: privateKey,
	}, nil
}

func createUnsignedTxnString(tx types.Transaction) (string, error) {
	encodedTx := msgpack.Encode(tx)

	sEnc := base64.StdEncoding.EncodeToString(encodedTx)

	return sEnc, nil

}

func compileProgram(client *algod.Client, programSource string) (compiledProgram []byte) {
	compileResponse, err := client.TealCompile([]byte(programSource)).Do(context.Background())
	if err != nil {
		fmt.Printf("Issue with compile: %s\n", err)
		return
	}
	compiledProgram, _ = base64.StdEncoding.DecodeString(compileResponse.Result)
	return compiledProgram
}
