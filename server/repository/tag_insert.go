package repository

import (
	"fmt"
	"strings"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func InsertTag(tag *UmiTag) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_tag (id, team_id, name, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		tag.Id, tag.TeamId, tag.Name, tag.CreatedAt, tag.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert tag: %w", err)
	}
	return nil
}

func BulkInsertTags(tags []*UmiTag) error {
	if len(tags) == 0 {
		return nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	rowPlaceholder := fmt.Sprintf("(%s)", utils.GenerateSqlPlaceholders(5))
	placeholders := make([]string, len(tags))
	args := make([]any, 0, len(tags)*5)
	for i, tag := range tags {
		placeholders[i] = rowPlaceholder
		args = append(args, tag.Id, tag.TeamId, tag.Name, tag.CreatedAt, tag.UpdatedAt)
	}
	query := fmt.Sprintf(`INSERT INTO umi_tag (id, team_id, name, created_at, updated_at) VALUES %s`, strings.Join(placeholders, ", "))
	_, err = connection.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to bulk insert tags: %w", err)
	}
	return nil
}
