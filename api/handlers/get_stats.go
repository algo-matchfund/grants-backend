package handlers

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
	"github.com/algo-matchfund/grants-backend/internal/database"
	"github.com/go-openapi/runtime/middleware"
)

type getStatsHandler handler

func (h *getStatsHandler) Handle(params operations.GetStatsParams) middleware.Responder {
	log.Printf("GET /stats by user")

	donationCount, err := h.db.GetDonationCount()
	if err != nil {
		log.Println(err)
		return operations.NewGetStatsInternalServerError()
	}

	donationAmount, err := h.db.GetTotalDonationAmount()
	if err != nil {
		log.Println(err)
		return operations.NewGetStatsInternalServerError()
	}

	matchingRound, err := h.db.GetCurrentMatchingRound()
	if err != nil {
		log.Println(err)
		return operations.NewGetStatsInternalServerError()
	}

	projectCount, err := h.db.GetProjectCount()
	if err != nil {
		log.Println(err)
		return operations.NewGetStatsInternalServerError()
	}

	userCount, err := h.db.GetNonCompanyUserCount()
	if err != nil {
		log.Println(err)
		return operations.NewGetStatsInternalServerError()
	}

	stats := new(models.Stats)
	stats.DonationAmount = &donationAmount
	stats.DonationCount = &donationCount
	stats.MatchAmount = &matchingRound.MatchAmount
	stats.MatchStartDate = matchingRound.StartDate
	stats.MatchEndDate = matchingRound.EndDate
	stats.ProjectCount = &projectCount
	stats.UserCount = &userCount

	return operations.NewGetStatsOK().WithPayload(stats)
}

// NewGetStatsHandler creates a handler for getting all stats
func NewGetStatsHandler(db *database.GrantsDatabase) operations.GetStatsHandler {
	return &getStatsHandler{db}
}
