package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func InsertEvent(event *UmiEvent, teamKey []byte) error {
	sqlEvent, err := NewSqlEvent(event, teamKey)
	if err != nil {
		return fmt.Errorf("failed to create sql event: %w", err)
	}
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_event (id, channel_id, title, body, level, action_url, icon_url, service_id, metadata, created_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		sqlEvent.Id, sqlEvent.ChannelId, sqlEvent.Title, sqlEvent.Body, sqlEvent.Level, sqlEvent.ActionURL, sqlEvent.IconURL,
		sqlEvent.ServiceId, sqlEvent.Metadata, sqlEvent.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event: %w", err)
	}
	return nil
}
