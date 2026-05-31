package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func CreateTeam(team *UmiTeam) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_team (id, name, encryption_key, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		team.Id, team.Name, team.EncryptionKey, team.CreatedAt, team.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create team: %w", err)
	}
	return nil
}
