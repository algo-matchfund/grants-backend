package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectsHandler handler

func (h *getProjectsHandler) Handle(params operations.GetProjectsParams) middleware.Responder {
	log.Printf("GET /projects")

	filter := &models.ProjectFilter{
		Name: params.Name,
	}

	projects, err := h.db.GetProjects(filter, params.Limit, params.Offset)

	if err != nil {
		log.Println(err)
		return operations.NewGetProjectsInternalServerError()
	}

	return operations.NewGetProjectsOK().WithPayload(projects)
}

// NewGetProjectsHandler creates a handler for getting all projects
func NewGetProjectsHandler(db *database.GrantsDatabase) operations.GetProjectsHandler {
	return &getProjectsHandler{db}
}
