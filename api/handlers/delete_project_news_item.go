package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postDeleteProjectsIDNewsNewsIDHandler handler

func (h *postDeleteProjectsIDNewsNewsIDHandler) Handle(params operations.DeleteProjectsIDNewsNewsIDParams, principal *models.Principal) middleware.Responder {
	return operations.NewDeleteProjectsIDNewsNewsIDOK()
}

// NewDeleteProjectsIDNewsNewsIDHandler creates a handler for creating a new project
func NewDeleteProjectsIDNewsNewsIDHandler(db *database.GrantsDatabase) operations.DeleteProjectsIDNewsNewsIDHandler {
	return &postDeleteProjectsIDNewsNewsIDHandler{db}
}
