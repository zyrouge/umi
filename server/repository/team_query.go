package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetTeamById(id string) (*UmiTeam, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, name, encryption_key, created_at, updated_at FROM umi_team WHERE id = ?`, id)
	team, err := SqlScanTeam(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return team, err
}

func ListTeamsByUserId(userId string) ([]*UmiTeam, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT t.id, t.name, t.encryption_key, t.created_at, t.updated_at FROM umi_team t
		 JOIN umi_member m ON m.team_id = t.id
		 WHERE m.user_id = ?
		 ORDER BY t.created_at ASC`,
		userId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list teams: %w", err)
	}
	defer rows.Close()
	var result []*UmiTeam
	for rows.Next() {
		team, err := SqlScanTeam(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, team)
	}
	return result, nil
}

func GetTeamByChannelId(channelId string) (*UmiTeam, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT t.id, t.name, t.encryption_key, t.created_at, t.updated_at FROM umi_team t
		 JOIN umi_channel c ON c.team_id = t.id
		 WHERE c.id = ?`,
		channelId,
	)
	team, err := SqlScanTeam(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return team, err
}
