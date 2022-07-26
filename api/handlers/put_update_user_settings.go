package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/go-openapi/runtime/middleware"
)

type putUsersSettingsHandler userHandler

func (h *putUsersSettingsHandler) Handle(params operations.PutUsersSettingsParams, principal *models.Principal) middleware.Responder {
	return operations.NewPutUsersSettingsOK()
}

// NewPutUsersSettingsHandler creates a handler for updating the authenticated user's settings
func NewPutUsersSettingsHandler(db *database.GrantsDatabase, keycloak *keycloak.KeycloakService) operations.PutUsersSettingsHandler {
	return &putUsersSettingsHandler{handler: handler{db}, keycloak: keycloak}
}
