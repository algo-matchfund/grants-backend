package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postNewProjectQuestionHandler handler

func (h *postNewProjectQuestionHandler) Handle(params operations.PostProjectQuestionParams, principal *models.Principal) middleware.Responder {
	return operations.NewPostProjectQuestionOK()
}

// NewPostProjectQuestionHandler creates a handler for creating a new question for a project
func NewPostProjectQuestionHandler(db *database.GrantsDatabase) operations.PostProjectQuestionHandler {
	return &postNewProjectQuestionHandler{db}
}
