package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type putUsersNotificationsHandler handler

func (h *putUsersNotificationsHandler) Handle(params operations.PutUsersNotificationsParams, principal *models.Principal) middleware.Responder {
	return operations.NewPutUsersNotificationsOK()
}

// NewPutUsersNotificationsHandler creates a handler for marking the authenticated user's notifications as read
func NewPutUsersNotificationsHandler(db *database.GrantsDatabase) operations.PutUsersNotificationsHandler {
	return &putUsersNotificationsHandler{db}
}
