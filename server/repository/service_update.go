package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func UpdateServiceName(id, name string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`UPDATE umi_service SET name = ?, updated_at = ? WHERE id = ?`, name, time.Now().Unix(), id)
	return err
}

func UpdateServiceTokenHash(id, tokenHash string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`UPDATE umi_service SET token_hash = ?, updated_at = ? WHERE id = ?`, tokenHash, time.Now().Unix(), id)
	return err
}
