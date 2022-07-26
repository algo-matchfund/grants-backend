package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectModerationByIDHandler handler

func (h *getProjectModerationByIDHandler) Handle(params operations.GetProjectModerationByIDParams, principal *models.Principal) middleware.Responder {
	log.Printf("GET /moderate/projects/%d", params.ModerationID)

	// if !principal.HasModeratorRole() {
	//   return operations.NewGetProjectModerationByIDForbidden()
	// }

	pendingProject, err := h.db.GetProjectModerationById(params.ModerationID)

	if err != nil {
		log.Println(err)
		return operations.NewGetProjectModerationByIDInternalServerError()
	}

	return operations.NewGetProjectModerationByIDOK().WithPayload(pendingProject)
}

// NewPostProjectsHandler creates a handler for getting a pending new project or project information change by moderation ID
func NewGetProjectModerationByIDHandler(db *database.GrantsDatabase) operations.GetProjectModerationByIDHandler {
	return &getProjectModerationByIDHandler{db}
}
