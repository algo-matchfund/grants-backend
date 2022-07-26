package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type putUsersNotificationsNotificationIDHandler handler

func (h *putUsersNotificationsNotificationIDHandler) Handle(params operations.PutUsersNotificationsNotificationIDParams, principal *models.Principal) middleware.Responder {
	return operations.NewPutUsersNotificationsNotificationIDOK()
}

// NewPutUsersNotificationsNotificationIDHandler creates a handler for marking the authenticated user's notification as read by its ID
func NewPutUsersNotificationsNotificationIDHandler(db *database.GrantsDatabase) operations.PutUsersNotificationsNotificationIDHandler {
	return &putUsersNotificationsNotificationIDHandler{db}
}
