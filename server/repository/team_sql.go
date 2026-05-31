package repository

import (
	"zyrouge.me/umi/utils"
)

func SqlScanTeam(scannable utils.SqlScannable) (*UmiTeam, error) {
	team := UmiTeam{}
	if err := scannable.Scan(&team.Id, &team.Name, &team.EncryptionKey, &team.CreatedAt, &team.UpdatedAt); err != nil {
		return nil, err
	}
	return &team, nil
}
