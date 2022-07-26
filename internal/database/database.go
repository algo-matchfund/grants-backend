package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	// Postgres driver for database/sql
	_ "github.com/lib/pq"

	// Needed for DB migration from files
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	sq "github.com/Masterminds/squirrel"
)

// Current version of the database, increment when adding new
// migration scripts
const migrationVersion = 1

// GrantsDatabase is a wrapper around sql.DB that adds methods to
// execute our analytics queries
type GrantsDatabase struct {
	*sql.DB

	// Used for building more complex queries with optional parameters
	// Defined here so that we can set the appropriate placeholder format
	// in one place (dollar sign for postgres)
	builder sq.StatementBuilderType

	driver string
}

// NewGrantsDatabase returns a new db
func NewGrantsDatabase(host, user, password, name string, port int) (*GrantsDatabase, error) {
	db := new(GrantsDatabase)
	db.driver = "postgres"
	db.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, name)

	var err error
	db.DB, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		return db, err
	}

	return db, nil
}

// Migrate migrates the DB schema to the latest version
func (db *GrantsDatabase) Migrate() error {
	// Retry until the database is available
	for err := db.Ping(); err != nil; err = db.Ping() {
		log.Printf("database unavailable, trying again in 5 seconds")
		time.Sleep(5 * time.Second)
	}

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})

	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migration",
		"postgres", driver)

	if err != nil {
		return err
	}

	v, dirty, _ := m.Version()

	if dirty {
		log.Fatalln("Database is in a dirty state, unable to migrate")
	}

	if v < migrationVersion {
		log.Printf("Migrating from migration %d to %d", v, migrationVersion)
		err = m.Migrate(migrationVersion)
	}

	return err
}

// newMockGrantsDatabase returns a new grants db
func newMockGrantsDatabase() (db *GrantsDatabase, mock sqlmock.Sqlmock) {
	db = new(GrantsDatabase)
	db.driver = "postgres"
	db.builder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	var err error
	db.DB, mock, err = sqlmock.New()

	if err != nil {
		panic(err)
	}

	return
}
