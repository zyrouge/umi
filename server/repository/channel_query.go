package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func GetChannelById(id string) (*UmiChannel, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(`SELECT id, name, team_id, created_at, updated_at FROM umi_channel WHERE id = ?`, id)
	ch, err := SqlScanChannel(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query channel by id: %w", err)
	}
	return ch, nil
}

func ListChannelsByTeamId(teamId string) ([]*UmiChannel, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT id, name, team_id, created_at, updated_at FROM umi_channel WHERE team_id = ? ORDER BY created_at ASC`,
		teamId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list channels by team id: %w", err)
	}
	defer rows.Close()
	var result []*UmiChannel
	for rows.Next() {
		ch, err := SqlScanChannel(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan channel by team id: %w", err)
		}
		result = append(result, ch)
	}
	return result, nil
}

func CountAccessibleChannelsByUserIdAndChannelIds(userId string, channelIds []string) (int, error) {
	if len(channelIds) == 0 {
		return 0, nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return 0, fmt.Errorf("failed to get database connection: %w", err)
	}
	query := fmt.Sprintf(
		`SELECT COUNT(DISTINCT c.id) FROM umi_channel c
		 JOIN member m ON m.team_id = c.team_id
		 WHERE m.user_id = ? AND c.id IN (%s)`,
		utils.GenerateSqlPlaceholders(len(channelIds)),
	)
	args := make([]any, 0, len(channelIds)+1)
	args = append(args, userId)
	args = append(args, utils.SliceToAny(channelIds)...)
	var count int
	if err = connection.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count accessible channels by user id and channel ids: %w", err)
	}
	return count, nil
}

func GetChannelTeamKeys(channelIds []string, masterKey []byte) (map[string][]byte, error) {
	if len(channelIds) == 0 {
		return map[string][]byte{}, nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	query := fmt.Sprintf(
		`SELECT c.id, t.encryption_key FROM umi_channel c
		 JOIN team t ON t.id = c.team_id
		 WHERE c.id IN (%s)`,
		utils.GenerateSqlPlaceholders(len(channelIds)),
	)
	args := utils.SliceToAny(channelIds)
	rows, err := connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query channel team keys by channel ids: %w", err)
	}
	defer rows.Close()
	result := make(map[string][]byte, len(channelIds))
	for rows.Next() {
		var channelId, encryptedKey string
		if err := rows.Scan(&channelId, &encryptedKey); err != nil {
			return nil, fmt.Errorf("failed to scan channel team key by channel id: %w", err)
		}
		decryptedKey, err := DecryptTeamKey(encryptedKey, masterKey)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt team key for channel %s: %w", channelId, err)
		}
		result[channelId] = decryptedKey
	}
	return result, nil
}
