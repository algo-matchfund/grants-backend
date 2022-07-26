package handlers

import (
	"fmt"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type postProjectsHandler handler

func (h *postProjectsHandler) Handle(params operations.PostProjectsParams, principal *models.Principal) middleware.Responder {
	log.Println("POST /projects")

	projectId, err := h.db.CreateProjectForModeration(params.Body, principal.ID)

	if err != nil {
		log.Println(err)
		return operations.NewPostProjectsInternalServerError()
	}

	newProject := new(models.Project)
	newProject.ID = fmt.Sprint(projectId)

	return operations.NewPostProjectsOK().WithPayload(newProject)
}

// NewPostProjectsHandler creates a handler for creating a new project
func NewPostProjectsHandler(db *database.GrantsDatabase) operations.PostProjectsHandler {
	return &postProjectsHandler{db}
}
