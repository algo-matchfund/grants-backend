package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectContributorsHandler handler

func (h *getProjectContributorsHandler) Handle(params operations.GetProjectContributorsParams) middleware.Responder {
	return operations.NewGetProjectContributorsOK()
}

// NewGetProjectContributorsHandler creates a handler for getting project's contributors
func NewGetProjectContributorsHandler(db *database.GrantsDatabase) operations.GetProjectContributorsHandler {
	return &getProjectContributorsHandler{db}
}
