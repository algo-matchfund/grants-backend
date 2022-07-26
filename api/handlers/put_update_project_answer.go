package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type updateAnswerHandler handler

func (h *updateAnswerHandler) Handle(params operations.UpdateAnswerParams, principal *models.Principal) middleware.Responder {
	return operations.NewUpdateAnswerOK()
}

// NewUpdateAnswerHandler creates a handler for updating an answer to a question about a project
func NewUpdateAnswerHandler(db *database.GrantsDatabase) operations.UpdateAnswerHandler {
	return &updateAnswerHandler{db}
}
