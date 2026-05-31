package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func DeleteExpiredEventsByRetentionDays(retentionDays int) error {
	if retentionDays <= 0 {
		return nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	maxCreatedAt := time.Now().UTC().AddDate(0, 0, -retentionDays)
	_, err = connection.Exec(`DELETE FROM umi_event WHERE created_at < ?`, maxCreatedAt.Unix())
	if err != nil {
		return fmt.Errorf("failed to delete expired events: %w", err)
	}
	return nil
}
