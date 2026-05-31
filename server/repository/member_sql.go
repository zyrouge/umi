package repository

import (
	"zyrouge.me/umi/utils"
)

func SqlScanMember(scannable utils.SqlScannable) (*UmiMember, error) {
	member := UmiMember{}
	if err := scannable.Scan(&member.UserId, &member.TeamId, &member.Role, &member.CreatedAt, &member.UpdatedAt); err != nil {
		return nil, err
	}
	return &member, nil
}
