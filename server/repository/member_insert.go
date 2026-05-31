package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func InsertMember(member *UmiMember) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_member (user_id, team_id, role, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		member.UserId, member.TeamId, member.Role, member.CreatedAt, member.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create member: %w", err)
	}
	return nil
}
