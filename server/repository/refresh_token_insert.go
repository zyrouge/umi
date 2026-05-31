package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func InsertRefreshToken(token *UmiRefreshToken) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_refresh_token (id, user_id, token_hash, expires_at, created_at) VALUES (?, ?, ?, ?, ?)`,
		token.Id, token.UserId, token.TokenHash, token.ExpiresAt, token.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert refresh token: %w", err)
	}
	return nil
}
