package repository

import (
	"zyrouge.me/umi/utils"
)

func SqlScanChannel(scannable utils.SqlScannable) (*UmiChannel, error) {
	channel := UmiChannel{}
	if err := scannable.Scan(&channel.Id, &channel.Name, &channel.TeamId, &channel.CreatedAt, &channel.UpdatedAt); err != nil {
		return nil, err
	}
	return &channel, nil
}
