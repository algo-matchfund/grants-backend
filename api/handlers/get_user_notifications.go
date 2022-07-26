package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getUsersNotificationsHandler handler

func (h *getUsersNotificationsHandler) Handle(params operations.GetUsersNotificationsParams, principal *models.Principal) middleware.Responder {
	return operations.NewGetUsersNotificationsOK()
}

// NewGetUsersNotificationsHandler creates a handler for getting the authenticated user's notifications
func NewGetUsersNotificationsHandler(db *database.GrantsDatabase) operations.GetUsersNotificationsHandler {
	return &getUsersNotificationsHandler{db}
}
