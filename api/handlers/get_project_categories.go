package handlers

import (
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getProjectCategoriesHandler handler

func (h *getProjectCategoriesHandler) Handle(params operations.GetProjectCategoriesParams) middleware.Responder {
	return operations.NewGetProjectCategoriesOK()
}

// NewGetProjectCategoriesHandler creates a handler for getting available project categories
func NewGetProjectCategoriesHandler(db *database.GrantsDatabase) operations.GetProjectCategoriesHandler {
	return &getProjectCategoriesHandler{db}
}
