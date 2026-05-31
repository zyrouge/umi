package repository

import (
	"fmt"

	"zyrouge.me/umi/database"
	"zyrouge.me/umi/utils"
)

func ListEventsByChannelIds(channelIds []string, limit int, channelKeys map[string][]byte) ([]*UmiEvent, error) {
	if len(channelIds) == 0 {
		return nil, nil
	}
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	query := fmt.Sprintf(
		`SELECT id, channel_id, title, body, level, action_url, icon_url, service_id, metadata, created_at
		 FROM umi_event WHERE channel_id IN (%s) ORDER BY created_at DESC LIMIT ?`,
		utils.GenerateSqlPlaceholders(len(channelIds)),
	)
	args := make([]any, 0, len(channelIds)+1)
	args = append(args, utils.SliceToAny(channelIds)...)
	args = append(args, limit)
	rows, err := connection.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query events: %w", err)
	}
	defer rows.Close()
	var result []*UmiEvent
	for rows.Next() {
		var raw UmiSqlEvent
		err := rows.Scan(
			&raw.Id, &raw.ChannelId, &raw.Title, &raw.Body, &raw.Level,
			&raw.ActionURL, &raw.IconURL, &raw.ServiceId, &raw.Metadata, &raw.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan raw event: %w", err)
		}
		teamKey, ok := channelKeys[raw.ChannelId]
		if !ok {
			return nil, fmt.Errorf("no team key for channel %s", raw.ChannelId)
		}
		event, err := raw.ToEvent(teamKey)
		if err != nil {
			return nil, fmt.Errorf("failed to convert sql event to event: %w", err)
		}
		result = append(result, event)
	}
	return result, nil
}
