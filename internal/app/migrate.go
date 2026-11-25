package app

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	// Migrate tools.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

const _defaultMigrationsPath = "file://migrations"

// Migrate runs database migrations.
func Migrate(pool *pgxpool.Pool) {
	log.Println("Running migrations...")

	connStr := pool.Config().ConnConfig.ConnString()

	m, err := migrate.New(_defaultMigrationsPath, connStr)
	if err != nil {
		log.Fatalf("migrate - New: %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatalf("migrate - Up: %v", err)
	}

	if errors.Is(err, migrate.ErrNoChange) {
		log.Println("No migrations to apply")
	} else {
		log.Println("Migrations applied successfully")
	}

	srcErr, dbErr := m.Close()
	if srcErr != nil {
		log.Printf("migrate - Close source: %v", srcErr)
	}
	if dbErr != nil {
		log.Printf("migrate - Close database: %v", dbErr)
	}

	fmt.Println()
}
