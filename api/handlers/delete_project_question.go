package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type deleteProjectQuestionHandler handler

func (h *deleteProjectQuestionHandler) Handle(params operations.DeleteQuestionParams, principal *models.Principal) middleware.Responder {
	return operations.NewDeleteQuestionOK()
}

// NewDeleteQuestionHandler creates a handler for deleting questions about a project
func NewDeleteQuestionHandler(db *database.GrantsDatabase) operations.DeleteQuestionHandler {
	return &deleteProjectQuestionHandler{db}
}
