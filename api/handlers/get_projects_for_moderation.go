package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectsForModerationHandler handler

func (h *getProjectsForModerationHandler) Handle(params operations.GetProjectsForModerationParams, principal *models.Principal) middleware.Responder {
	log.Println("GET /moderate/projects")

	// if !principal.HasModeratorRole() {
	//   return operations.NewGetProjectsForModerationForbidden()
	// }

	pendingProjects, err := h.db.GetProjectsForModeration(params.Name, params.Limit, params.Offset)

	if err != nil {
		log.Println(err)
		return operations.NewGetProjectsForModerationInternalServerError()
	}

	return operations.NewGetProjectsForModerationOK().WithPayload(pendingProjects)
}

// NewPostProjectsHandler creates a handler for getting list of pending new projects and project changes
func NewGetProjectsForModerationHandler(db *database.GrantsDatabase) operations.GetProjectsForModerationHandler {
	return &getProjectsForModerationHandler{db}
}
