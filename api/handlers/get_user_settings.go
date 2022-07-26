package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/go-openapi/runtime/middleware"
)

type getUsersSettingsHandler userHandler

func (h *getUsersSettingsHandler) Handle(params operations.GetUsersSettingsParams, principal *models.Principal) middleware.Responder {
	return operations.NewGetUsersSettingsOK()
}

// NewGetUsersSettingsHandler creates a handler for getting the authenticated user's settings
func NewGetUsersSettingsHandler(db *database.GrantsDatabase, keycloak *keycloak.KeycloakService) operations.GetUsersSettingsHandler {
	return &getUsersSettingsHandler{handler: handler{db}, keycloak: keycloak}
}
