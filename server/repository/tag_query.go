package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func GetTagById(id string) (*UmiTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT id, team_id, name, created_at, updated_at FROM umi_tag WHERE id = ?`, id,
	)
	tag, err := SqlScanTag(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tag by id: %w", err)
	}
	return tag, nil
}

func GetTagByTeamIdAndName(teamId string, name string) (*UmiTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT id, team_id, name, created_at, updated_at FROM umi_tag WHERE team_id = ? AND name = ?`,
		teamId, name,
	)
	tag, err := SqlScanTag(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query tag by team id and name: %w", err)
	}
	return tag, nil
}

func GetTagByNames(teamId string, names []string) (map[string]*UmiTag, error) {
	if len(names) == 0 {
		return nil, nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	result := make(map[string]*UmiTag)
	query := fmt.Sprintf(
		`SELECT id, team_id, name, created_at, updated_at FROM umi_tag WHERE team_id = ? AND name IN (%s)`,
		utils.GenerateSqlPlaceholders(len(names)),
	)
	args := append([]any{teamId}, utils.SliceToAny(names)...)
	rows, err := connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query tags by names: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		tag, err := SqlScanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan tag: %w", err)
		}
		result[tag.Name] = tag
	}
	return result, nil
}

func ListTagsByTeamId(teamId string) ([]*UmiTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT id, team_id, name, created_at, updated_at FROM umi_tag WHERE team_id = ? ORDER BY created_at ASC`,
		teamId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}
	defer rows.Close()
	var result []*UmiTag
	for rows.Next() {
		tag, err := SqlScanTag(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, tag)
	}
	return result, nil
}
