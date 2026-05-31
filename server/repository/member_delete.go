package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func DeleteMemberByUserIdAndTeamId(userId string, teamId string) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(`DELETE FROM umi_member WHERE user_id = ? AND team_id = ?`, userId, teamId)
	if err != nil {
		return fmt.Errorf("failed to delete member: %w", err)
	}
	return nil
}
