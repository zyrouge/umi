package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func UpdateChannelNameById(id string, name string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	now := time.Now()
	_, err = connection.Exec(`UPDATE umi_channel SET name = ?, updated_at = ? WHERE id = ?`, name, now.Unix(), id)
	return err
}
