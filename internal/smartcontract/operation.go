package smartcontract

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/algo-matchfund/grants-backend/internal/config"
	"github.com/go-openapi/strfmt"

	"github.com/algorand/go-algorand-sdk/client/v2/algod"
	"github.com/algorand/go-algorand-sdk/crypto"
	"github.com/algorand/go-algorand-sdk/future"
	"github.com/algorand/go-algorand-sdk/types"
)

type SmartContractClient struct {
	client                  *algod.Client
	blockConfirmationWindow uint64
	creatorMnemonic         string
}

func NewSmartContractClient(config *config.Config) (*SmartContractClient, error) {
	algodClient, err := algod.MakeClient(config.SmartContract.Node.Address, config.SmartContract.Node.Token)
	if err != nil {
		return nil, err
	}

	sc := SmartContractClient{
		client:                  algodClient,
		blockConfirmationWindow: config.SmartContract.BlockConfirmation,
		creatorMnemonic:         config.SmartContract.Admin.Passphrase,
	}

	if !sc.init() {
		return nil, errors.New("algod client failed health check")
	}

	return &sc, nil
}

func (s *SmartContractClient) init() bool {
	// up to 10 seconds for health check
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := s.client.HealthCheck().Do(ctx)
	if err != nil {
		log.Printf("SmartContractClient.Init: Algod client failed health check, error - %s\n", err)
		return false
	}

	return true
}

func (s *SmartContractClient) CreateApp(endDate strfmt.DateTime) (int64, error) {
	// get account from mnemonic
	creatorAccount, err := getAccount(s.creatorMnemonic)
	if err != nil {
		return 0, err
	}

	// get transaction suggested parameters
	params, err := s.client.SuggestedParams().Do(context.Background())
	if err != nil {
		return 0, err
	}

	// declare application state storage (immutatble)
	localInts := 1
	localBytes := 0
	globalInts := 6
	globalBytes := 1
	globalSchema := types.StateSchema{NumUint: uint64(globalInts), NumByteSlice: uint64(globalBytes)}
	localSchema := types.StateSchema{NumUint: uint64(localInts), NumByteSlice: uint64(localBytes)}

	// read smart contract sources
	a, err := ioutil.ReadFile("approval.teal")
	if err != nil {
		fmt.Printf("Cannot read the approval program: %s\n", err)
		return 0, err
	}
	c, err := ioutil.ReadFile("clear.teal")
	if err != nil {
		fmt.Printf("Cannot read the clear program: %s\n", err)
		return 0, err
	}

	approvalProgram := compileProgram(s.client, string(a[:]))
	clearProgram := compileProgram(s.client, string(c[:]))

	appArgs := make([][]byte, 3)
	// to timestamp
	now := uint64(time.Now().Unix())
	end := uint64(time.Time(endDate).Unix())
	nowByte := make([]byte, 8)
	binary.BigEndian.PutUint64(nowByte, now)
	endByte := make([]byte, 8)
	binary.BigEndian.PutUint64(endByte, end)

	// creator address
	appArgs[0] = []byte(creatorAccount.Address.String())
	// start time
	appArgs[1] = nowByte
	// end time
	appArgs[2] = endByte

	// create unsigned transaction
	txn, err := future.MakeApplicationCreateTxWithExtraPages(false, approvalProgram, clearProgram, globalSchema, localSchema,
		appArgs, nil, nil, nil, params, creatorAccount.Address, nil,
		types.Digest{}, [32]byte{}, types.Address{}, 1)
	if err != nil {
		return 0, err
	}

	// Sign the transaction
	txID, signedTxn, err := crypto.SignTransaction(creatorAccount.PrivateKey, txn)
	fmt.Printf("Signed txid: %s\n", txID)
	if err != nil {
		return 0, err
	}

	// Submit the transaction
	sendResponse, err := s.client.SendRawTransaction(signedTxn).Do(context.Background())
	if err != nil {
		return 0, err
	}
	fmt.Printf("Submitted transaction %s\n", sendResponse)

	// Wait for confirmation
	confirmedTxn, err := future.WaitForConfirmation(s.client, txID, 4, context.Background())
	if err != nil {
		fmt.Printf("Error waiting for confirmation on txID: %s\n", txID)
		return 0, err
	}
	fmt.Printf("Confirmed Transaction: %s in Round %d\n", txID, confirmedTxn.ConfirmedRound)

	// display results
	confirmedTxn, _, err = s.client.PendingTransactionInformation(txID).Do(context.Background())
	if err != nil {
		return 0, err
	}
	appId := confirmedTxn.ApplicationIndex
	fmt.Printf("Created new app-id: %d\n", appId)

	return int64(appId), nil
}

func (s *SmartContractClient) CreateUnsignedOptInTxn(appId uint64, senderAdress string) (string, error) {
	params, err := s.client.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	// turns input address string to address object
	address, err := types.DecodeAddress(senderAdress)
	if err != nil {
		return "", err
	}

	// create unsigned transaction
	tx, err := future.MakeApplicationOptInTx(appId, nil, nil, nil, nil, params,
		address, nil, types.Digest{}, [32]byte{}, types.Address{})
	if err != nil {
		return "", err
	}

	txString, err := createUnsignedTxnString(tx)
	if err != nil {
		return "", err
	}

	return txString, nil
}

func (s *SmartContractClient) CreateSetDonation(appId uint64, senderAdress string, amount int) (string, error) {
	if amount <= 0 {
		return "", errors.New("amount must be bigger than 0")
	}

	params, err := s.client.SuggestedParams().Do(context.Background())
	if err != nil {
		return "", err
	}

	appArgs := make([][]byte, 2)
	appArgs[0] = []byte("set_donation")
	appArgs[1] = make([]byte, 8)
	binary.BigEndian.PutUint64(appArgs[1], uint64(amount))

	// turns input address string to address object
	address, err := types.DecodeAddress(senderAdress)
	if err != nil {
		return "", err
	}

	// create unsigned transaction
	tx, err := future.MakeApplicationNoOpTx(appId, appArgs, nil, nil, nil, params, address,
		nil, types.Digest{}, [32]byte{}, types.Address{})

	if err != nil {
		return "", err
	}

	txString, err := createUnsignedTxnString(tx)
	if err != nil {
		return "", err
	}

	return txString, nil

}
