package repository

import (
	"fmt"
	"time"

	"zyrouge.me/umi/database"
)

func UpdateMemberRole(userId string, teamId string, role UmiMemberRole) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`UPDATE umi_member SET role = ?, updated_at = ? WHERE user_id = ? AND team_id = ?`,
		role, time.Now().Unix(), userId, teamId,
	)
	return err
}
