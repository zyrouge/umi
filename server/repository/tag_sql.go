package repository

import (
	"fmt"

	"zyrouge.me/umi/utils"
)

func SqlScanTag(scannable utils.SqlScannable) (*UmiTag, error) {
	tag := UmiTag{}
	if err := scannable.Scan(&tag.Id, &tag.TeamId, &tag.Name, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan tag: %w", err)
	}
	return &tag, nil
}
