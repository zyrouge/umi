package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
)

func CreateService(service *UmiService) error {
	connection, err := database.GetConnection()
	if err != nil {
		return fmt.Errorf("failed to get database connection: %w", err)
	}
	_, err = connection.Exec(
		`INSERT INTO umi_service (id, name, team_id, token_hash, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		service.Id, service.Name, service.TeamId, service.TokenHash, service.CreatedAt, service.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create service: %w", err)
	}
	return nil
}
