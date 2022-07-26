package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectNewsHandler handler

func (h *getProjectNewsHandler) Handle(params operations.GetProjectNewsParams) middleware.Responder {
	return operations.NewGetProjectNewsOK()
}

// NewGetProjectNewsHandler creates a handler for getting project's news
func NewGetProjectNewsHandler(db *database.GrantsDatabase) operations.GetProjectNewsHandler {
	return &getProjectNewsHandler{db}
}
