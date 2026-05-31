package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetRefreshTokenById(id string) (*UmiRefreshToken, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, created_at FROM umi_refresh_token WHERE id = ?`,
		id,
	)
	token, err := SqlScanRefreshToken(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return token, err
}

func GetRefreshTokenByTokenHash(tokenHash string) (*UmiRefreshToken, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT id, user_id, token_hash, expires_at, created_at FROM umi_refresh_token WHERE token_hash = ?`,
		tokenHash,
	)
	token, err := SqlScanRefreshToken(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return token, err
}

func ListRefreshTokensByUserId(userId string) ([]*UmiRefreshToken, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT id, user_id, token_hash, expires_at, created_at FROM umi_refresh_token WHERE user_id = ? ORDER BY created_at ASC`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list refresh tokens: %w", err)
	}
	defer rows.Close()
	var result []*UmiRefreshToken
	for rows.Next() {
		token, err := SqlScanRefreshToken(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, token)
	}
	return result, nil
}
