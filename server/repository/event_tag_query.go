package repository

import (
	"database/sql"
	"errors"
	"fmt"

	"zyrouge.me/umi/database"
)

func GetEventTagAssociation(eventId string, tagId string) (*UmiEventTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	row := connection.QueryRow(
		`SELECT event_id, tag_id, created_at FROM umi_event_tag_map WHERE event_id = ? AND tag_id = ?`,
		eventId, tagId,
	)
	eventTag, err := SqlScanEventTag(row)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query event tag association by event id and tag id: %w", err)
	}
	return eventTag, nil
}

func ListEventTagAssociationsByEventId(eventId string) ([]*UmiEventTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT event_id, tag_id, created_at FROM umi_event_tag_map WHERE event_id = ? ORDER BY created_at ASC`,
		eventId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list event tag associations by event id: %w", err)
	}
	defer rows.Close()
	var result []*UmiEventTag
	for rows.Next() {
		eventTag, err := SqlScanEventTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event tag association by event id: %w", err)
		}
		result = append(result, eventTag)
	}
	return result, nil
}

func ListEventTagsByEventId(eventId string) ([]*UmiTag, error) {
	connection, err := database.GetConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	rows, err := connection.Query(
		`SELECT t.id, t.team_id, t.name, t.created_at, t.updated_at
		 FROM umi_tag t
		 JOIN umi_event_tag_map et ON et.tag_id = t.id
		 WHERE et.event_id = ?
		 ORDER BY t.name ASC`,
		eventId,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list event tags by event id: %w", err)
	}
	defer rows.Close()
	var result []*UmiTag
	for rows.Next() {
		tag, err := SqlScanTag(rows)
		if err != nil {
			return nil, fmt.Errorf("failed to scan event tag by event id: %w", err)
		}
		result = append(result, tag)
	}
	return result, nil
}
