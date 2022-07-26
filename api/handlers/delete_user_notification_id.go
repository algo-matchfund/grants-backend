package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type deleteUsersNotificationsNotificationIDHandler handler

func (h *deleteUsersNotificationsNotificationIDHandler) Handle(params operations.DeleteUsersNotificationsNotificationIDParams, principal *models.Principal) middleware.Responder {
	return operations.NewDeleteUsersNotificationsNotificationIDOK()
}

// NewDeleteUsersNotificationsNotificationIDHandler creates a handler for deleteting the authenticated user's notifications
func NewDeleteUsersNotificationsNotificationIDHandler(db *database.GrantsDatabase) operations.DeleteUsersNotificationsNotificationIDHandler {
	return &deleteUsersNotificationsNotificationIDHandler{db}
}
