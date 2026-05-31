package repository

import (
	"fmt"
	"strings"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func InsertEventTag(eventTag *UmiEventTag) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_event_tag_map (event_id, tag_id, created_at) VALUES (?, ?, ?)`,
		eventTag.EventId, eventTag.TagId, eventTag.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert event tag: %w", err)
	}
	return nil
}

func BulkInsertEventTags(eventTags []*UmiEventTag) error {
	if len(eventTags) == 0 {
		return nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	rowPlaceholder := fmt.Sprintf("(%s)", utils.GenerateSqlPlaceholders(3))
	placeholders := make([]string, len(eventTags))
	args := make([]any, 0, len(eventTags)*3)
	for i, et := range eventTags {
		placeholders[i] = rowPlaceholder
		args = append(args, et.EventId, et.TagId, et.CreatedAt)
	}
	query := fmt.Sprintf(`INSERT INTO umi_event_tag_map (event_id, tag_id, created_at) VALUES %s`, strings.Join(placeholders, ", "))
	_, err = connection.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert event tags: %w", err)
	}
	return nil
}
