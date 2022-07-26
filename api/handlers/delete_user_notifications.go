package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type deleteUsersNotificationsHandler handler

func (h *deleteUsersNotificationsHandler) Handle(params operations.DeleteUsersNotificationsParams, principal *models.Principal) middleware.Responder {
	return operations.NewDeleteUsersNotificationsOK()
}

// NewDeleteUsersNotificationsHandler creates a handler for deleteting the authenticated user's notification by ID
func NewDeleteUsersNotificationsHandler(db *database.GrantsDatabase) operations.DeleteUsersNotificationsHandler {
	return &deleteUsersNotificationsHandler{db}
}
