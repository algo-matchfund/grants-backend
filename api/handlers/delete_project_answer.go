package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type deleteAnswerHandler handler

func (h *deleteAnswerHandler) Handle(params operations.DeleteAnswerParams, principal *models.Principal) middleware.Responder {
	return operations.NewDeleteAnswerOK()
}

// NewDeleteAnswerHandler creates a handler for deleting an answer to a question about a project
func NewDeleteAnswerHandler(db *database.GrantsDatabase) operations.DeleteAnswerHandler {
	return &deleteAnswerHandler{db}
}
