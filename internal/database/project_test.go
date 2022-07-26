package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

const (
	projectName = "test"
)

func Test_GetProjectCount(t *testing.T) {
	db, mock := newMockGrantsDatabase()

	mock.ExpectQuery(`SELECT COUNT(.*) FROM projects`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).FromCSVString("20"))

	count, err := db.GetProjectCount()
	if err != nil {
		t.Error(err)
	}

	if count != 20 {
		t.Errorf("Expected 20 companies, got %d\n", count)
	}
}
func Test_GetProjects(t *testing.T) {
	db, mock := newMockGrantsDatabase()

	rows := sqlmock.NewRows([]string{"name", "icon", "description", "id"}).
		AddRow(projectName, "", "", 1)

	mock.ExpectQuery(`
		SELECT name, icon, description, id
		FROM projects
		ORDER BY name asc`).
		WithArgs(1).
		WillReturnRows(rows)

	// //var i int64 = 1
	// companies, err := db.GetProjects(nil, nil, nil)

	// assertStrEq(*companies[0].Name, projectName, t)

	// if err != nil {
	// 	t.Error(err)
	// }

	// if err = mock.ExpectationsWereMet(); err != nil {
	// 	t.Errorf("there were unfulfilled expectations: %s", err)
	// }

	defer db.Close()
}
