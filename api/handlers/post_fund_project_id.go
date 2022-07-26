package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postNewProjectIDFundHandler handler

func (h *postNewProjectIDFundHandler) Handle(params operations.PostProjectIDFundParams, principal *models.Principal) middleware.Responder {
	log.Printf("POST /project/%s/fund", params.ID)

	_, err := h.db.PostFund(params.ID, params.Amount, principal.ID)

	if err != nil {
		log.Println(err)
		return operations.NewPostProjectIDFundInternalServerError()
	}

	return operations.NewPostProjectIDFundOK()
}

// NewPostProjectIDFundHandler creates a handler for funding a project
func NewPostProjectIDFundHandler(db *database.GrantsDatabase) operations.PostProjectIDFundHandler {
	return &postNewProjectIDFundHandler{db}
}
