package migration

import (
	"database/sql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"log"
)

func MigrateUp(connectionString string) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/recon?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///scripts",
		"postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
}
