package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func UpdateTeamName(id, name string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`UPDATE umi_team SET name = ?, updated_at = ? WHERE id = ?`, name, time.Now().Unix(), id)
	return err
}
