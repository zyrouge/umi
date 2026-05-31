package migrations

import (
	"database/sql"
	"fmt"
	"time"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

type Migration struct {
	Version int
	Up      func(transaction *sql.Tx) error
	Down    func(transaction *sql.Tx) error
}

var Migrations []Migration = []Migration{
	MigrationV1,
}

func MigrateDatabase() error {
	utils.Logger.Info().Msg("starting database migrations...")
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`
		CREATE TABLE IF NOT EXISTS umi_migration (
			version    INTEGER NOT NULL,
			applied_at BIGINT NOT NULL,
			CONSTRAINT pk_umi_migration PRIMARY KEY (version)
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	for _, x := range Migrations {
		var count int
		if err := connection.QueryRow(`SELECT COUNT(*) FROM umi_migration WHERE version = ?`, x.Version).Scan(&count); err != nil {
			return fmt.Errorf("failed to check migration %d: %w", x.Version, err)
		}
		if count > 0 {
			utils.Logger.Info().Int("version", x.Version).Msg("migration already applied, skipping...")
			continue
		}
		if err := PerformMigration(connection, x); err != nil {
			return err
		}
	}
	return nil
}

func PerformMigration(connection *sql.DB, migration Migration) error {
	utils.Logger.Info().Int("version", migration.Version).Msg("applying migration...")
	transaction, err := connection.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction for migration %d: %w", migration.Version, err)
	}
	defer transaction.Rollback()
	if err := migration.Up(transaction); err != nil {
		return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
	}
	now := time.Now()
	if _, err := transaction.Exec(`INSERT INTO umi_migration (version, applied_at) VALUES (?, ?)`, migration.Version, now.Unix()); err != nil {
		return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
	}
	if err := transaction.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
	}
	utils.Logger.Info().Int("version", migration.Version).Msg("migration applied successfully")
	return nil
}
