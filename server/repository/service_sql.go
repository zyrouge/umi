package repository

import (
	"zyrouge.me/umi/utils"
)

func SqlScanService(scannable utils.SqlScannable) (*UmiService, error) {
	service := UmiService{}
	if err := scannable.Scan(&service.Id, &service.Name, &service.TeamId, &service.TokenHash, &service.CreatedAt, &service.UpdatedAt); err != nil {
		return nil, err
	}
	return &service, nil
}
