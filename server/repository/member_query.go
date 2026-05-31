package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetMember(userId string, teamId string) (*UmiMember, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT user_id, team_id, role, created_at, updated_at FROM umi_member WHERE user_id = ? AND team_id = ?`,
		userId, teamId,
	)
	member, err := SqlScanMember(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return member, err
}

func ListMembersByTeamId(teamId string) ([]*UmiMember, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT user_id, team_id, role, created_at, updated_at FROM umi_member WHERE team_id = ? ORDER BY created_at ASC`,
		teamId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list members: %w", err)
	}
	defer rows.Close()
	var result []*UmiMember
	for rows.Next() {
		member, err := SqlScanMember(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, member)
	}
	return result, nil
}
