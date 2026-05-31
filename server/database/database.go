package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	_ "modernc.org/sqlite"
	"zyrouge.me/umi/application"
	"zyrouge.me/umi/utils"
)

var cachedConnection *sql.DB

func LoadConnection() error {
	config, err := application.GetConfig()
	if err != nil {
		return err
	}
	driver := string(config.Database.Driver)
	utils.Logger.Info().Str("driver", driver).Msg("connecting to database")
	connection, err := sql.Open(driver, config.Database.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}
	if err := connection.Ping(); err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	connection.SetMaxOpenConns(config.Database.MaxOpenConnections)
	connection.SetMaxIdleConns(config.Database.MaxIdleConnections)
	cachedConnection = connection
	utils.Logger.Info().Str("driver", driver).Msg("database connected")
	return nil
}

func GetConnection() (*sql.DB, error) {
	if cachedConnection == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}
	return cachedConnection, nil
}
