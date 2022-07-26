package database

import (
	"github.com/algo-matchfund/grants-backend/gen/models"
)

func (db *GrantsDatabase) GetCurrentMatchingRound() (matchingRound models.MatchingRound, err error) {
	err = db.QueryRow(`
		SELECT id, start_date, end_date, match_amount
		FROM matching_rounds
		ORDER BY id desc limit 1`).
		Scan(&matchingRound.ID, &matchingRound.StartDate, &matchingRound.EndDate, &matchingRound.MatchAmount)

	return
}
