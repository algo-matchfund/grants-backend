package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/algo-matchfund/grants-backend/gen/models"
	"github.com/go-openapi/strfmt"
)

// Currently only supports moderation for new projects, but not project changes.
func (db *GrantsDatabase) PostProjectModerationById(moderationId, appId int64, body models.PendingProjectModeration, moderatorId string) error {
	var pendingProjectId int64

	tx, err := db.Begin()

	if err != nil {
		tx.Rollback()
		return err
	}

	stmt, params, err := db.builder.
		Update("moderations").
		Set("status", body.Status).
		Set("comment", body.Comment).
		Set("moderator_id", moderatorId).
		Where("id = ?", moderationId).
		Suffix("returning after_project_id").
		ToSql()

	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.QueryRow(stmt, params...).Scan(&pendingProjectId)

	if err != nil {
		tx.Rollback()
		return err
	}

	if body.Status == "approve" {
		var projectId int64

		err = db.QueryRow(`
		INSERT INTO projects (name, description, algorand_wallet, created_at, created_by, icon, background, image, content, app_id)
		SELECT name, description, algorand_wallet, created_at, created_by, icon, background, image, content, $1
		FROM pending_projects
		WHERE id = $2
		RETURNING id
		`, appId, pendingProjectId).Scan(&projectId)

		if err != nil {
			tx.Rollback()
			return err
		}

		_, err = db.Exec(`
		INSERT INTO project_socials (project_id, github, twitter, email, website, other)
		SELECT $1, github, twitter, email, website, other
		FROM pending_projects
		WHERE id = $2
		`, projectId, pendingProjectId)

		if err != nil {
			tx.Rollback()
			return err
		}

		// notify user that project changes have been approved
	} else {
		// logic for "deny", probably involves notification.
	}

	err = tx.Commit()

	return err
}

func (db *GrantsDatabase) GetProjectsForModeration(name *string, limit *int64, offset *int64) ([]*models.PendingProject, error) {
	pendingProjects := []*models.PendingProject{}

	query := db.builder.
		Select(`m.id, m.created_at, m.before_project_id, m.after_project_id,
		p1.name, p1.description, p1.created_at, p1.icon, p1.background, p1.image, p1.content,
		p1.github, p1.twitter, p1.email, p1.website,
		p2.name, p2.description, p2.created_at, p2.icon, p2.background, p2.image, p2.content,
		p2.github, p2.twitter, p2.email, p2.website
		`).
		From("moderations m").
		LeftJoin("pending_projects p1 on p1.id = m.before_project_id").
		LeftJoin("pending_projects p2 on p2.id = m.after_project_id").
		Where("m.moderator_id is null") // Get pending projects that are not assigned to any moderator

	if limit != nil {
		query = query.Limit(uint64(*limit))
	}

	if offset != nil {
		query = query.Offset(uint64(*offset))
	}

	if name != nil {
		query = query.Where("p1.name like ? or p2.name like ?", fmt.Sprintf("%%%s%%", *name), fmt.Sprintf("%%%s%%", *name))
	}

	stmt, params := query.MustSql()
	rows, err := db.Query(stmt, params...)

	if err != nil {
		return pendingProjects, err
	}

	for rows.Next() {
		pendingProject := new(models.PendingProject)
		pendingProject.Before = new(models.Project)
		pendingProject.Before.Socials = new(models.Socials)
		pendingProject.After = new(models.Project)
		pendingProject.After.Socials = new(models.Socials)

		var beforeId sql.NullString
		var afterId sql.NullString
		var beforeName sql.NullString
		var beforeDescription sql.NullString
		var beforeCreatedAt sql.NullTime
		var beforeIcon sql.NullString
		var beforeBackground sql.NullString
		var beforeImage sql.NullString
		var beforeContent sql.NullString
		var beforeGithub sql.NullString
		var beforeTwitter sql.NullString
		var beforeEmail sql.NullString
		var beforeWebsite sql.NullString
		var afterName sql.NullString
		var afterDescription sql.NullString
		var afterCreatedAt sql.NullTime
		var afterIcon sql.NullString
		var afterBackground sql.NullString
		var afterImage sql.NullString
		var afterContent sql.NullString
		var afterGithub sql.NullString
		var afterTwitter sql.NullString
		var afterEmail sql.NullString
		var afterWebsite sql.NullString

		err = rows.Scan(&pendingProject.ModerationID,
			&pendingProject.Date,
			&beforeId,
			&afterId,
			&beforeName,
			&beforeDescription,
			&beforeCreatedAt,
			&beforeIcon,
			&beforeBackground,
			&beforeImage,
			&beforeContent,
			&beforeGithub,
			&beforeTwitter,
			&beforeEmail,
			&beforeWebsite,
			&afterName,
			&afterDescription,
			&afterCreatedAt,
			&afterIcon,
			&afterBackground,
			&afterImage,
			&afterContent,
			&afterGithub,
			&afterTwitter,
			&afterEmail,
			&afterWebsite)

		if err != nil {
			log.Println(err)
			continue
		}

		// Project before modification must have a name. Without a name would mean no "before project"
		if beforeName.Valid {
			pendingProject.Before.ID = beforeId.String
			pendingProject.Before.Name = &beforeName.String
			pendingProject.Before.Description = beforeDescription.String
			pendingProject.Before.CreatedAt = strfmt.DateTime(beforeCreatedAt.Time)
		}

		if beforeIcon.Valid {
			pendingProject.Before.Icon = beforeIcon.String
		}

		if beforeBackground.Valid {
			pendingProject.Before.Background = beforeBackground.String
		}

		if beforeImage.Valid {
			pendingProject.Before.Image = beforeImage.String
		}

		if beforeContent.Valid {
			pendingProject.Before.Content = beforeContent.String
		}

		if beforeGithub.Valid {
			pendingProject.Before.Socials.Github = beforeGithub.String
		}

		if beforeTwitter.Valid {
			pendingProject.Before.Socials.Twitter = beforeTwitter.String
		}

		if beforeEmail.Valid {
			pendingProject.Before.Socials.Email = beforeEmail.String
		}

		if beforeWebsite.Valid {
			pendingProject.Before.Socials.Web = beforeWebsite.String
		}

		if afterId.Valid {
			pendingProject.After.ID = afterId.String
		}

		if afterName.Valid {
			pendingProject.After.Name = &afterName.String
		}

		if afterDescription.Valid {
			pendingProject.After.Description = afterDescription.String
		}

		if afterCreatedAt.Valid {
			pendingProject.After.CreatedAt = strfmt.DateTime(afterCreatedAt.Time)
		}

		if afterIcon.Valid {
			pendingProject.After.Icon = afterIcon.String
		}

		if afterBackground.Valid {
			pendingProject.After.Background = afterBackground.String
		}

		if afterImage.Valid {
			pendingProject.After.Image = afterImage.String
		}

		if afterContent.Valid {
			pendingProject.After.Content = afterContent.String
		}

		if afterGithub.Valid {
			pendingProject.After.Socials.Github = afterGithub.String
		}

		if afterTwitter.Valid {
			pendingProject.After.Socials.Twitter = afterTwitter.String
		}

		if afterEmail.Valid {
			pendingProject.After.Socials.Email = afterEmail.String
		}

		if afterWebsite.Valid {
			pendingProject.After.Socials.Web = afterWebsite.String
		}

		pendingProjects = append(pendingProjects, pendingProject)
	}

	return pendingProjects, nil
}

func (db *GrantsDatabase) GetProjectModerationById(moderationId int64) (*models.PendingProject, error) {
	pendingProject := new(models.PendingProject)
	pendingProject.Before = new(models.Project)
	pendingProject.Before.Socials = new(models.Socials)
	pendingProject.After = new(models.Project)
	pendingProject.After.Socials = new(models.Socials)

	query := db.builder.
		Select(`m.id, m.created_at, m.before_project_id, m.after_project_id,
		p1.name, p1.description, p1.created_at, p1.icon, p1.background, p1.image, p1.content,
		p1.github, p1.twitter, p1.email, p1.website,
		p2.name, p2.description, p2.created_at, p2.icon, p2.background, p2.image, p2.content,
		p2.github, p2.twitter, p2.email, p2.website
		`).
		From("moderations m").
		LeftJoin("pending_projects p1 on p1.id = m.before_project_id").
		LeftJoin("pending_projects p2 on p2.id = m.after_project_id").
		Where("m.id = ?", moderationId)

	stmt, params := query.MustSql()
	row := db.QueryRow(stmt, params...)

	if row.Err() != nil {
		return nil, row.Err()
	}

	var beforeId sql.NullString
	var afterId sql.NullString
	var beforeName sql.NullString
	var beforeDescription sql.NullString
	var beforeCreatedAt sql.NullTime
	var beforeIcon sql.NullString
	var beforeBackground sql.NullString
	var beforeImage sql.NullString
	var beforeContent sql.NullString
	var beforeGithub sql.NullString
	var beforeTwitter sql.NullString
	var beforeEmail sql.NullString
	var beforeWebsite sql.NullString
	var afterName sql.NullString
	var afterDescription sql.NullString
	var afterCreatedAt sql.NullTime
	var afterIcon sql.NullString
	var afterBackground sql.NullString
	var afterImage sql.NullString
	var afterContent sql.NullString
	var afterGithub sql.NullString
	var afterTwitter sql.NullString
	var afterEmail sql.NullString
	var afterWebsite sql.NullString

	err := row.Scan(&pendingProject.ModerationID,
		&pendingProject.Date,
		&beforeId,
		&afterId,
		&beforeName,
		&beforeDescription,
		&beforeCreatedAt,
		&beforeIcon,
		&beforeBackground,
		&beforeImage,
		&beforeContent,
		&beforeGithub,
		&beforeTwitter,
		&beforeEmail,
		&beforeWebsite,
		&afterName,
		&afterDescription,
		&afterCreatedAt,
		&afterIcon,
		&afterBackground,
		&afterImage,
		&afterContent,
		&afterGithub,
		&afterTwitter,
		&afterEmail,
		&afterWebsite)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Project before modification must have a name. Without a name would mean no "before project"
	if beforeName.Valid {
		pendingProject.Before.ID = beforeId.String
		pendingProject.Before.Name = &beforeName.String
		pendingProject.Before.Description = beforeDescription.String
		pendingProject.Before.CreatedAt = strfmt.DateTime(beforeCreatedAt.Time)
	}

	if beforeIcon.Valid {
		pendingProject.Before.Icon = beforeIcon.String
	}

	if beforeBackground.Valid {
		pendingProject.Before.Background = beforeBackground.String
	}

	if beforeImage.Valid {
		pendingProject.Before.Image = beforeImage.String
	}

	if beforeContent.Valid {
		pendingProject.Before.Content = beforeContent.String
	}

	if beforeGithub.Valid {
		pendingProject.Before.Socials.Github = beforeGithub.String
	}

	if beforeTwitter.Valid {
		pendingProject.Before.Socials.Twitter = beforeTwitter.String
	}

	if beforeEmail.Valid {
		pendingProject.Before.Socials.Email = beforeEmail.String
	}

	if beforeWebsite.Valid {
		pendingProject.Before.Socials.Web = beforeWebsite.String
	}

	if afterId.Valid {
		pendingProject.After.ID = afterId.String
	}

	if afterName.Valid {
		pendingProject.After.Name = &afterName.String
	}

	if afterDescription.Valid {
		pendingProject.After.Description = afterDescription.String
	}

	if afterCreatedAt.Valid {
		pendingProject.After.CreatedAt = strfmt.DateTime(afterCreatedAt.Time)
	}

	if afterIcon.Valid {
		pendingProject.After.Icon = afterIcon.String
	}

	if afterBackground.Valid {
		pendingProject.After.Background = afterBackground.String
	}

	if afterImage.Valid {
		pendingProject.After.Image = afterImage.String
	}

	if afterContent.Valid {
		pendingProject.After.Content = afterContent.String
	}

	if afterGithub.Valid {
		pendingProject.After.Socials.Github = afterGithub.String
	}

	if afterTwitter.Valid {
		pendingProject.After.Socials.Twitter = afterTwitter.String
	}

	if afterEmail.Valid {
		pendingProject.After.Socials.Email = afterEmail.String
	}

	if afterWebsite.Valid {
		pendingProject.After.Socials.Web = afterWebsite.String
	}

	return pendingProject, nil
}
