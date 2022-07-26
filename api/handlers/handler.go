package handlers

import (
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
)

type handler struct {
	db *database.GrantsDatabase
}

type userHandler struct {
	handler
	keycloak *keycloak.KeycloakService
}
