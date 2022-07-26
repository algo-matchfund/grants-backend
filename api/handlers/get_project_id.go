package handlers

import (
	"database/sql"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectsIdHandler handler

func (h *getProjectsIdHandler) Handle(params operations.GetProjectsIDParams) middleware.Responder {
	log.Printf("GET /projects/%s", params.ID)

	project, err := h.db.GetProjectById(params.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return operations.NewGetProjectsIDNotFound()
		}

		log.Println(err)
		return operations.NewGetProjectsIDInternalServerError()
	}

	return operations.NewGetProjectsIDOK().WithPayload(project)
}

// NewGetProjectsIDHandler creates a handler for getting a project by ID
func NewGetProjectsIDHandler(db *database.GrantsDatabase) operations.GetProjectsIDHandler {
	return &getProjectsIdHandler{db}
}
