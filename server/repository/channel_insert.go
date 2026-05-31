package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func CreateChannel(channel *UmiChannel) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_channel (id, name, team_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		channel.Id, channel.Name, channel.TeamId, channel.CreatedAt, channel.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create channel: %w", err)
	}
	return nil
}
