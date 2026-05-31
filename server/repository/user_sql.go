package repository

import (
	"zyrouge.me/umi/utils"
)

func SqlScanUser(scannable utils.SqlScannable) (*UmiUser, error) {
	user := UmiUser{}
	if err := scannable.Scan(&user.Id, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt, &user.UpdatedAt); err != nil {
		return nil, err
	}
	return &user, nil
}
