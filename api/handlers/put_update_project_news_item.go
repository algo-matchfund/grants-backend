package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type putProjectNewsItemHandler handler

func (h *putProjectNewsItemHandler) Handle(params operations.UpdateProjectNewsItemParams, principal *models.Principal) middleware.Responder {
	return operations.NewUpdateProjectNewsItemOK()
}

// NewUpdateProjectNewsItemHandler creates a handler for updating a single project news item
func NewUpdateProjectNewsItemHandler(db *database.GrantsDatabase) operations.UpdateProjectNewsItemHandler {
	return &putProjectNewsItemHandler{db}
}
