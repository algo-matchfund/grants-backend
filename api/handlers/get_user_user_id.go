package handlers

import (
	"database/sql"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/keycloak"
	"github.com/go-openapi/runtime/middleware"
)

type getUsersUserIDHandler userHandler

func (h *getUsersUserIDHandler) Handle(params operations.GetUsersUserIDParams) middleware.Responder {
	log.Printf("GET /user/%s", params.UserID)

	// if !principal.HasAdminRole() {
	//   return operations.NewGetUserUserIDForbidden()
	// }

	user, err := h.db.GetUser(params.UserID)

	if err != nil {
		if err == sql.ErrNoRows {
			return operations.NewGetUsersUserIDNotFound()
		}

		log.Println(err)
		return operations.NewGetUsersUserIDInternalServerError()
	}

	return operations.NewGetUsersUserIDOK().WithPayload(user)
}

func NewGetUserUserIdHandler(db *database.GrantsDatabase, keycloak *keycloak.KeycloakService) operations.GetUsersUserIDHandler {
	return &getUsersUserIDHandler{handler: handler{db}, keycloak: keycloak}
}
