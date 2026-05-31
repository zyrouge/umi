package repository

import (
	"fmt"

	"zyrouge.me/umi/utils"
)

func SqlScanEventTag(scannable utils.SqlScannable) (*UmiEventTag, error) {
	eventTag := UmiEventTag{}
	if err := scannable.Scan(&eventTag.EventId, &eventTag.TagId, &eventTag.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan event tag: %w", err)
	}
	return &eventTag, nil
}
