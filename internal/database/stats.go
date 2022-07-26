package database

import (
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
)

// GetStats returns all statistics the provided user has read permission on
func (db *GrantsDatabase) GetStats() (models.Stats, error) {
	stats := models.Stats{}
	query := db.builder.
		Select(`donationCount, matchAmount, projectCount`).
		From("stats")

	stmt, params := query.MustSql()
	rows, err := db.Query(stmt, params...)

	if err != nil {
		return stats, err
	}

	for rows.Next() {
		// stats = new(models.Stats)
		err = rows.Scan(&stats.DonationCount, &stats.MatchAmount, &stats.ProjectCount)

		if err != nil {
			log.Println(err)
			continue
		}
	}

	return stats, nil
}
