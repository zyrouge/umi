package application

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
	"zyrouge.me/umi/utils"
)

type UmiConfig struct {
	Server   UmiServerConfig   `toml:"server" validate:"required"`
	Database UmiDatabaseConfig `toml:"database" validate:"required"`
	Secret   UmiSecretConfig   `toml:"secret" validate:"required"`
}

type UmiServerConfig struct {
	Host        string   `toml:"host" validate:"required"`
	Port        int      `toml:"port" validate:"required"`
	CrossOrigin []string `toml:"cross_origin"`
	WebFiles    string   `toml:"web_files"`
}

type UmiDatabaseConfigDriver string

const (
	UmiDatabaseConfigDriverSQLite   UmiDatabaseConfigDriver = "sqlite"
	UmiDatabaseConfigDriverPostgres UmiDatabaseConfigDriver = "postgres"
)

type UmiDatabaseConfig struct {
	Driver             UmiDatabaseConfigDriver `toml:"driver" validate:"required,oneof=sqlite postgres"`
	DSN                string                  `toml:"dsn" validate:"required"`
	RetentionDays      int                     `toml:"retention_days"`
	MaxOpenConnections int                     `toml:"max_open_connections"`
	MaxIdleConnections int                     `toml:"max_idle_connections"`
}

type UmiSecretConfig struct {
	JwtSecret         string `toml:"jwt_secret" validate:"required"`
	TeamEncryptionKey string `toml:"team_encryption_key" validate:"required"`
	UserEncryptionKey string `toml:"user_encryption_key" validate:"required"`

	JwtSecretBytes         []byte
	TeamEncryptionKeyBytes []byte
	UserEncryptionKeyBytes []byte
}

var configCache *UmiConfig

func LoadConfig() error {
	file := os.Getenv("CONFIG_FILE")
	if file == "" {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		file = filepath.Join(cwd, "config.toml")
	}
	utils.Logger.Info().Str("file", file).Msg("loading config...")
	bytes, err := os.ReadFile(file)
	if err != nil {
		return err
	}
	var config UmiConfig
	if err = toml.Unmarshal(bytes, &config); err != nil {
		return err
	}
	if err := utils.GlobalValidator.Struct(&config); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}
	jwtSecretBytes, err := base64.StdEncoding.DecodeString(config.Secret.JwtSecret)
	if err != nil {
		return fmt.Errorf("failed to decode jwt secret: %w", err)
	}
	teamEncryptionKeyBytes, err := base64.StdEncoding.DecodeString(config.Secret.TeamEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decode team encryption key: %w", err)
	}
	userEncryptionKeyBytes, err := base64.StdEncoding.DecodeString(config.Secret.UserEncryptionKey)
	if err != nil {
		return fmt.Errorf("failed to decode user encryption key: %w", err)
	}
	config.Secret.JwtSecretBytes = jwtSecretBytes
	config.Secret.TeamEncryptionKeyBytes = teamEncryptionKeyBytes
	config.Secret.UserEncryptionKeyBytes = userEncryptionKeyBytes
	configCache = &config
	utils.Logger.Info().Str("file", file).Msg("config loaded")
	return nil
}

func GetConfig() (*UmiConfig, error) {
	if configCache == nil {
		return nil, fmt.Errorf("config is not loaded")
	}
	return configCache, nil
}
