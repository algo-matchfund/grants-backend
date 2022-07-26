package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/go-openapi/runtime/middleware"
)

type getUsersHandler userHandler

func (h *getUsersHandler) Handle(params operations.GetUsersParams, principal *models.Principal) middleware.Responder {
	return operations.NewGetUsersOK()
}

// NewGetUsersHandler creates a handler for getting the authenticated user
func NewGetUsersHandler(db *database.GrantsDatabase, keycloak *keycloak.KeycloakService) operations.GetUsersHandler {
	return &getUsersHandler{handler: handler{db}, keycloak: keycloak}
}
