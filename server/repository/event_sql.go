package repository

import (
	"encoding/json"
	"fmt"

	"zyrouge.me/umi/utils"
)

type UmiSqlEvent struct {
	Id        string
	ChannelId string
	Title     string
	Body      *string
	Level     *string
	ActionURL *string
	IconURL   *string
	ServiceId string
	Metadata  *string
	CreatedAt int64
}

func NewSqlEvent(event *UmiEvent, teamKey []byte) (*UmiSqlEvent, error) {
	titleEncrypted, err := utils.EncryptAESGCM(teamKey, event.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt event title: %w", err)
	}
	var bodyEncrypted *string
	if event.Body != nil {
		v, err := utils.EncryptAESGCM(teamKey, *event.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event body: %w", err)
		}
		bodyEncrypted = &v
	}
	var levelEncrypted *string
	if event.Level != nil {
		v, err := utils.EncryptAESGCM(teamKey, string(*event.Level))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event level: %w", err)
		}
		levelEncrypted = &v
	}
	var actionURLEncrypted *string
	if event.ActionURL != nil {
		v, err := utils.EncryptAESGCM(teamKey, *event.ActionURL)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event action_url: %w", err)
		}
		actionURLEncrypted = &v
	}
	var iconURLEncrypted *string
	if event.IconURL != nil {
		v, err := utils.EncryptAESGCM(teamKey, *event.IconURL)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event icon_url: %w", err)
		}
		iconURLEncrypted = &v
	}
	var metadataEncrypted *string
	if event.Metadata != nil {
		metaJson, err := json.Marshal(event.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal event metadata: %w", err)
		}
		v, err := utils.EncryptAESGCM(teamKey, string(metaJson))
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt event metadata: %w", err)
		}
		metadataEncrypted = &v
	}
	sqlEvent := UmiSqlEvent{
		Id:        event.Id,
		ChannelId: event.ChannelId,
		Title:     titleEncrypted,
		Body:      bodyEncrypted,
		Level:     levelEncrypted,
		ActionURL: actionURLEncrypted,
		IconURL:   iconURLEncrypted,
		ServiceId: event.ServiceId,
		Metadata:  metadataEncrypted,
		CreatedAt: event.CreatedAt,
	}
	return &sqlEvent, nil
}

func (sqlEvent *UmiSqlEvent) ToEvent(teamKey []byte) (*UmiEvent, error) {
	event := UmiEvent{
		Id:        sqlEvent.Id,
		ChannelId: sqlEvent.ChannelId,
		ServiceId: sqlEvent.ServiceId,
		CreatedAt: sqlEvent.CreatedAt,
	}
	title, err := utils.DecryptAESGCM(teamKey, sqlEvent.Title)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt event title: %w", err)
	}
	event.Title = title
	if sqlEvent.Body != nil {
		v, err := utils.DecryptAESGCM(teamKey, *sqlEvent.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event body: %w", err)
		}
		event.Body = &v
	}
	if sqlEvent.Level != nil {
		v, err := utils.DecryptAESGCM(teamKey, *sqlEvent.Level)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event level: %w", err)
		}
		level := UmiEventLevel(v)
		event.Level = &level
	}
	if sqlEvent.ActionURL != nil {
		v, err := utils.DecryptAESGCM(teamKey, *sqlEvent.ActionURL)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event action_url: %w", err)
		}
		event.ActionURL = &v
	}
	if sqlEvent.IconURL != nil {
		v, err := utils.DecryptAESGCM(teamKey, *sqlEvent.IconURL)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event icon_url: %w", err)
		}
		event.IconURL = &v
	}
	if sqlEvent.Metadata != nil {
		v, err := utils.DecryptAESGCM(teamKey, *sqlEvent.Metadata)
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt event metadata: %w", err)
		}
		meta := map[string]string{}
		if err := json.Unmarshal([]byte(v), &meta); err != nil {
			return nil, fmt.Errorf("failed to unmarshal event metadata: %w", err)
		}
		event.Metadata = meta
	}
	return &event, nil
}

func SqlScanEvent(scannable utils.SqlScannable, teamKey []byte) (*UmiEvent, error) {
	var sqlEvent UmiSqlEvent
	err := scannable.Scan(
		&sqlEvent.Id, &sqlEvent.ChannelId, &sqlEvent.Title, &sqlEvent.Body, &sqlEvent.Level,
		&sqlEvent.ActionURL, &sqlEvent.IconURL, &sqlEvent.ServiceId, &sqlEvent.Metadata, &sqlEvent.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return sqlEvent.ToEvent(teamKey)
}
