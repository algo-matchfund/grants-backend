package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postNewProjectNewsItemHandler handler

func (h *postNewProjectNewsItemHandler) Handle(params operations.PostProjectNewsParams, principal *models.Principal) middleware.Responder {
	return operations.NewPostProjectNewsOK()
}

// NewPostProjectNewsHandler creates a handler for creating a new project news item
func NewPostProjectNewsHandler(db *database.GrantsDatabase) operations.PostProjectNewsHandler {
	return &postNewProjectNewsItemHandler{db}
}
