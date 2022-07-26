package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postProjectAnswerHandler handler

func (h *postProjectAnswerHandler) Handle(params operations.PostProjectAnswerParams, principal *models.Principal) middleware.Responder {
	return operations.NewPostProjectAnswerOK()
}

// NewPostProjectAnswerHandler creates a handler for submitting a new answer to a question about a project
func NewPostProjectAnswerHandler(db *database.GrantsDatabase) operations.PostProjectAnswerHandler {
	return &postProjectAnswerHandler{db}
}
