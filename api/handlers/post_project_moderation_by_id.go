package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/algo-matchfund/grants-backend/internal/smartcontract"
	"github.com/go-openapi/runtime/middleware"
)

type postProjectModerationByIDHandler struct {
	handler
	smartContractClient *smartcontract.SmartContractClient
}

func (h *postProjectModerationByIDHandler) Handle(params operations.PostProjectModerationByIDParams, principal *models.Principal) middleware.Responder {
	log.Printf("POST /moderate/projects/%d", params.ModerationID)

	// if !principal.HasModeratorRole() {
	//   return operations.NewPostProjectModerationByIDForbidden()
	// }

	var appId int64

	if params.Body.Status == "approve" {
		matchingRound, err := h.db.GetCurrentMatchingRound()
		if err != nil {
			log.Println(err)
			return operations.NewGetStatsInternalServerError()
		}

		appId, err = h.smartContractClient.CreateApp(matchingRound.EndDate)
		if err != nil {
			log.Println(err)
			return operations.NewPostProjectModerationByIDInternalServerError()
		}
	}

	err := h.db.PostProjectModerationById(params.ModerationID, appId, *params.Body, principal.ID)
	if err != nil {
		log.Println(err)
		return operations.NewPostProjectModerationByIDInternalServerError()
	}

	return operations.NewPostProjectModerationByIDOK()
}

// NewPostProjectsHandler creates a handler for submitting moderation for a project by moderation ID
func NewPostProjectModerationByIDHandler(db *database.GrantsDatabase, scc *smartcontract.SmartContractClient) operations.PostProjectModerationByIDHandler {
	return &postProjectModerationByIDHandler{handler: handler{db}, smartContractClient: scc}
}
