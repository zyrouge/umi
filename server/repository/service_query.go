package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetServiceById(id string) (*UmiService, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, name, team_id, token_hash, created_at, updated_at FROM umi_service WHERE id = ?`, id)
	service, err := SqlScanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return service, err
}

func GetServiceByTokenHash(tokenHash string) (*UmiService, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, name, team_id, token_hash, created_at, updated_at FROM umi_service WHERE token_hash = ?`, tokenHash)
	service, err := SqlScanService(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	return service, err
}

func ListServicesByTeamId(teamId string) ([]*UmiService, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT id, name, team_id, token_hash, created_at, updated_at FROM umi_service WHERE team_id = ? ORDER BY created_at ASC`,
		teamId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()
	var result []*UmiService
	for rows.Next() {
		service, err := SqlScanService(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, service)
	}
	return result, nil
}
