package handlers

import (
	"log"
	"strconv"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/smartcontract"
	"github.com/go-openapi/runtime/middleware"
)

type getSmartContractTransactionsHandler struct {
	smartContractClient *smartcontract.SmartContractClient
}

func (h *getSmartContractTransactionsHandler) Handle(params operations.GetSmartContractTransactionsParams, principal *models.Principal) middleware.Responder {
	log.Printf("GET /transactions/%s", params.ID)

	id, err := (strconv.Atoi(params.ID))
	if err != nil {
		log.Println(err)
		return operations.NewGetSmartContractTransactionsBadRequest()
	}
	appId := uint64(id)

	optinTxn, err := h.smartContractClient.CreateUnsignedOptInTxn(appId, params.Address)

	if err != nil {
		log.Println(err)
		return operations.NewGetSmartContractTransactionsInternalServerError()
	}

	setDonationTxn, err := h.smartContractClient.CreateSetDonation(appId, params.Address, int(params.Amount))

	if err != nil {
		log.Println(err)
		return operations.NewGetSmartContractTransactionsInternalServerError()
	}

	pl := models.Transactions{Optin: &optinTxn, SetDonation: &setDonationTxn}

	return operations.NewGetSmartContractTransactionsOK().WithPayload(&pl)
}

// NewGetSmartContractHandler creates a handler for getting unsigned transactions by app ID
func NewGetSmartContractHandler(scc *smartcontract.SmartContractClient) operations.GetSmartContractTransactionsHandler {
	return &getSmartContractTransactionsHandler{smartContractClient: scc}
}
