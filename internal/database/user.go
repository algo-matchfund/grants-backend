package database

import (
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/algo-matchfund/grants-backend/gen/models"
)

func (db *GrantsDatabase) GetUsers(limit, offset *int64) ([]*models.User, error) {
	var query squirrel.SelectBuilder
	query = db.builder.
		Select(`
      id,
      first_name,
      last_name,
      email`).
		From("users")

	if limit != nil {
		query = query.Limit(uint64(*limit))
	}

	if offset != nil {
		query = query.Offset(uint64(*offset))
	}

	stmt, params := query.MustSql()

	rows, err := db.Query(stmt, params...)

	if err != nil {
		return nil, err
	}

	users := []*models.User{}

	for rows.Next() {
		user := new(models.User)
		firstName := new(sql.NullString)
		lastName := new(sql.NullString)
		err := rows.Scan(&user.ID, firstName, lastName, &user.Email)

		if firstName.Valid {
			user.FirstName = firstName.String
		}

		if lastName.Valid {
			user.LastName = lastName.String
		}

		if err != nil {
			log.Println(err)
			continue
		}

		users = append(users, user)
	}

	for _, user := range users {
		user.Roles, err = db.GetUserRoles(user.ID)

		if err != nil {
			log.Println(err)
		}
	}

	return users, nil
}

func (db *GrantsDatabase) GetUser(id string) (*models.User, error) {
	stmt, params := db.builder.
		Select(`
		id,
		first_name,
		last_name,
		email
		`).
		From("users").
		Where("id = ?", id).MustSql()

	row := db.QueryRow(stmt, params...)

	if row.Err() != nil {
		return nil, row.Err()
	}

	user := new(models.User)
	firstName := new(sql.NullString)
	lastName := new(sql.NullString)

	err := row.Scan(&user.ID, firstName, lastName, &user.Email)
	if err != nil {
		return nil, err
	}

	if firstName.Valid {
		user.FirstName = firstName.String
	}

	if lastName.Valid {
		user.LastName = lastName.String
	}

	user.Roles, err = db.GetUserRoles(user.ID)

	return user, err
}

func (db *GrantsDatabase) CreateUser(userId string) (string, error) {
	stmt, params := db.builder.
		Insert("users").
		Columns("id").
		Values(userId).
		Suffix("returning id").MustSql()

	var id string
	err := db.QueryRow(stmt, params...).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (db *GrantsDatabase) GetUserRoles(userId string) ([]string, error) {
	roles := []string{}

	stmt, params := db.builder.
		Select("role").
		From("user_roles").
		Where("user_id = ?", userId).MustSql()

	rows, err := db.Query(stmt, params...)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var role string
		err = rows.Scan(&role)

		if err != nil {
			log.Println(err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func (db *GrantsDatabase) GetNonCompanyUserCount() (int64, error) {
	var count int64

	err := db.QueryRow("SELECT COUNT(*) FROM users where is_company IS NOT TRUE").Scan(&count)

	return count, err
}
