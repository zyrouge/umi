package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func DeleteRefreshToken(id string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`DELETE FROM umi_refresh_token WHERE id = ?`, id)
	return err
}

func DeleteRefreshTokensByUserId(userId string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`DELETE FROM umi_refresh_token WHERE user_id = ?`, userId)
	return err
}
