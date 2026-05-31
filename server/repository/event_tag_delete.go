package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func DeleteEventTagAssociationsByEventIds(eventIds []string) error {
	if len(eventIds) == 0 {
		return nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	query := fmt.Sprintf(`DELETE FROM umi_event_tag_map WHERE event_id IN (%s)`, utils.GenerateSqlPlaceholders(len(eventIds)))
	args := utils.SliceToAny(eventIds)
	_, err = connection.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete event tag associations by event ids: %w", err)
	}
	return nil
}
