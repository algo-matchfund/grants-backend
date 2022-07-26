package database

import (
	"database/sql"
	"fmt"
	"log"

	// "github.com/lib/pq"
	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/algo-matchfund/grants-backend/gen/restapi/operations"
)

// getProjectDonors returns the total number of contributors in a project
func (db *GrantsDatabase) getProjectDonors(projectId string) (int64, error) {
	var count int64

	stmt, params := db.builder.
		Select("count(distinct user_id)").
		From("contributions").
		Where("project_id = ?", projectId).
		MustSql()

	row := db.QueryRow(stmt, params...)

	if row.Err() != nil {
		return count, row.Err()
	}

	err := row.Scan(&count)

	return count, err
}

// GetProjectCount returns the total number of projects in the database
func (db *GrantsDatabase) GetProjectCount() (int64, error) {
	var count int64

	err := db.QueryRow("SELECT COUNT(*) FROM projects").Scan(&count)

	return count, err
}

func (db *GrantsDatabase) GetProjects(filter *models.ProjectFilter, limit *int64, offset *int64) ([]*models.Project, error) {
	projects := []*models.Project{}
	query := db.builder.
		Select(`p.name, p.algorand_wallet, p.icon, p.description, p.id, p.created_at, p.app_id ,m.match`).
		From("projects p").
		LeftJoin("matches m on p.id = m.project_id").
		OrderBy("p.name asc")

	if limit != nil {
		query = query.Limit(uint64(*limit))
	}

	if offset != nil {
		query = query.Offset(uint64(*offset))
	}

	if filter != nil && filter.Name != nil {
		query = query.Where("p.name like ?", fmt.Sprintf("%%%s%%", *filter.Name))
	}

	stmt, params := query.MustSql()
	rows, err := db.Query(stmt, params...)

	if err != nil {
		return projects, err
	}

	for rows.Next() {
		project := new(models.Project)
		var fund Fund
		var projectMatch sql.NullFloat64

		err = rows.Scan(&project.Name, &project.AlgorandWallet, &project.Icon, &project.Description, &project.ID, &project.CreatedAt, &project.AppID, &projectMatch)

		if err != nil {
			log.Println(err)
			continue
		}
		fund.ProjectId = project.ID

		// Get contributions
		contributions := []*models.ProjectContributor{}
		fundAmount := int64(0)
		q := db.builder.
			Select("amount, user_id, id").
			From("contributions").
			Where("project_id = ?", project.ID)

		s, p := q.MustSql()
		r, err := db.Query(s, p...)

		if err != nil {
			log.Println(err)
			continue
		}

		if projectMatch.Valid {
			project.Match = projectMatch.Float64
		} else {
			project.Match = 0.0
		}

		for r.Next() {
			contribution := new(models.ProjectContributor)
			err = r.Scan(&contribution.Amount, &contribution.ContributorID, &contribution.ID)

			if err != nil {
				log.Println(err)
				continue
			}

			fund.Amount = append(fund.Amount, uint64(contribution.Amount))
			fundAmount += contribution.Amount
			contributions = append(contributions, contribution)
		}

		project.Contributions = contributions
		project.FundAmount = fundAmount

		donors, err := db.getProjectDonors(project.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		project.Donors = donors
		projects = append(projects, project)
	}

	return projects, nil
}

func (db *GrantsDatabase) GetProjectById(id string) (*models.Project, error) {
	project := new(models.Project)

	query := db.builder.
		Select(`p.name, p.algorand_wallet, p.icon, p.image, p.description, p.id, p.content, p.created_at, p.app_id, m.match`).
		From("projects p").
		LeftJoin("matches m on p.id = m.project_id").
		Where("p.id = ?", id)

	stmt, params := query.MustSql()
	row := db.QueryRow(stmt, params...)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var projectMatch sql.NullFloat64
	var image sql.NullString

	err := row.Scan(&project.Name, &project.AlgorandWallet, &project.Icon, &image, &project.Description, &project.ID, &project.Content, &project.CreatedAt, &project.AppID, &projectMatch)
	if err != nil {
		return nil, err
	}

	if projectMatch.Valid {
		project.Match = projectMatch.Float64
	} else {
		project.Match = 0.0
	}

	if image.Valid {
		project.Image = image.String
	}

	contributions := []*models.ProjectContributor{}
	fundAmount := int64(0)
	q := db.builder.
		Select("amount, user_id, id").
		From("contributions").
		Where("project_id = ?", project.ID)
	s, p := q.MustSql()
	r, err := db.Query(s, p...)

	if err != nil {
		log.Println(err)
		return project, err
	}

	for r.Next() {
		contribution := new(models.ProjectContributor)
		err = r.Scan(&contribution.Amount, &contribution.ContributorID, &contribution.ID)

		if err != nil {
			log.Println(err)
			continue
		}

		fundAmount += contribution.Amount
		contributions = append(contributions, contribution)
	}

	project.Contributions = contributions
	project.FundAmount = fundAmount

	donors, err := db.getProjectDonors(project.ID)
	if err != nil {
		return project, err
	}
	project.Donors = donors

	return project, err
}

func (db *GrantsDatabase) GetAlgorandWalletByProjectId(projectId string) (string, error) {
	query := db.builder.
		Select(`algorand_wallet`).
		From("projects").
		Where("id = ?", projectId)

	stmt, params := query.MustSql()
	row := db.QueryRow(stmt, params...)

	if row.Err() != nil {
		return "", row.Err()
	}

	var wallet string
	err := row.Scan(&wallet)
	if err != nil {
		return "", err
	}

	return wallet, nil
}

func (db *GrantsDatabase) CreateProjectForModeration(body operations.PostProjectsBody, userId string) (int64, error) {
	var pendingProjectId int64

	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		return pendingProjectId, err
	}

	stmt, params, err := db.builder.
		Insert("pending_projects").
		Columns("name", "algorand_wallet", "description", "icon", "image", "content", "created_by", "github", "website").
		Values(body.Name, body.AlgorandWallet, body.ShortDescription, body.Icon, body.Screenshot, body.Content, userId, body.Github, body.Homepage).
		Suffix("returning id").ToSql()

	if err != nil {
		tx.Rollback()
		return pendingProjectId, err
	}

	err = tx.QueryRow(stmt, params...).Scan(&pendingProjectId)
	if err != nil {
		tx.Rollback()
		return pendingProjectId, err
	}

	stmt, params, err = db.builder.
		Insert("moderations").
		Columns("after_project_id").
		Values(pendingProjectId).
		ToSql()

	if err != nil {
		tx.Rollback()
		return pendingProjectId, err
	}

	_, err = tx.Exec(stmt, params...)
	if err != nil {
		tx.Rollback()
		return pendingProjectId, err
	}

	err = tx.Commit()

	return pendingProjectId, err
}
