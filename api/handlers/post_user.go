package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/go-openapi/runtime/middleware"
)

type postUserHandler userHandler

func (h *postUserHandler) Handle(params operations.PostUsersParams, principal *models.Principal) middleware.Responder {
	log.Printf("POST /users")

	_, err := h.db.CreateUser(params.UserID)

	if err != nil {
		log.Println(err)
		return operations.NewPostUsersInternalServerError()
	}

	return operations.NewPostUsersOK()
}

// NewPostUserHandler creates a handler for registering a new user account
func NewPostUserHandler(db *database.GrantsDatabase, keycloak *keycloak.KeycloakService) operations.PostUsersHandler {
	return &postUserHandler{handler: handler{db}, keycloak: keycloak}
}
