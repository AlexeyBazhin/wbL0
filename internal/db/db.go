package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

var (
	//go:embed migrations/*.sql
	fs embed.FS
)

// ConnectDatabase creates instance of db connection
func ConnectDatabase(ctx context.Context, driverName string, dsn string) (*sql.DB, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return db, fmt.Errorf("соединение с базой данных не сформированно: %w", err)
	}

	db.SetMaxOpenConns(2)
	db.SetMaxIdleConns(10)

	for i := 0; i < 10; i++ {
		if err = db.Ping(); err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		break
	}
	if err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	if err = migrateUp(dsn); err != nil {
		return nil, fmt.Errorf("failed to migrate db: %w", err)
	}

	return db, nil
}

// ConnectPostgreSQL form instance of PostgresQL connection
func ConnectPostgreSQL(ctx context.Context, dsn string) (*sql.DB, error) {
	return ConnectDatabase(ctx, "postgres", dsn)
}

func migrateUp(url string) error {
	sourceInstance, err := iofs.New(fs, "migrations")
	if err != nil {
		return err
	}
	m, err := migrate.NewWithSourceInstance("iofs", sourceInstance, url)
	if err != nil {
		return fmt.Errorf("failed to create new migrate instance: %w", err)
	}

	if err = m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations up: %w", err)
	}

	return nil
}
