package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectQAHandler handler

func (h *getProjectQAHandler) Handle(params operations.GetProjectQAParams) middleware.Responder {
	return operations.NewGetProjectQAOK()
}

// NewGetProjectQAHandler creates a handler for getting project's questions and answers
func NewGetProjectQAHandler(db *database.GrantsDatabase) operations.GetProjectQAHandler {
	return &getProjectQAHandler{db}
}
