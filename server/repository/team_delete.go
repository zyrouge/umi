package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func DeleteTeam(id string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`DELETE FROM umi_team WHERE id = ?`, id)
	return err
}
