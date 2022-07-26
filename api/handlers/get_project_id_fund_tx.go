package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/service/watchdog"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectIDFundTxHandler struct {
	handler
	watchdogFactory *watchdog.WatchdogFactory
}

func (h *getProjectIDFundTxHandler) Handle(params operations.GetProjectIDFundTxParams, principal *models.Principal) middleware.Responder {
	log.Printf("GET /project/%s/fund/tx", params.ID)

	txs, err := h.db.GetAlgorandTransactions(params.ID, principal.ID, *params.Pending)
	if err != nil {
		return operations.NewGetProjectIDFundTxInternalServerError()
	}

	return operations.NewGetProjectIDFundTxOK().WithPayload(txs)
}

// NewgetProjectIDFundTxHandler creates a handler for getting project's contributors
func NewGetProjectIDFundTxHandler(db *database.GrantsDatabase, watchdogFactory *watchdog.WatchdogFactory) operations.GetProjectIDFundTxHandler {
	return &getProjectIDFundTxHandler{
		handler:         handler{db},
		watchdogFactory: watchdogFactory,
	}
}
