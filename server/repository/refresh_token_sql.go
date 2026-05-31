package repository

import (
	"fmt"

	"zyrouge.me/umi/utils"
)

func SqlScanRefreshToken(scannable utils.SqlScannable) (*UmiRefreshToken, error) {
	token := UmiRefreshToken{}
	if err := scannable.Scan(&token.Id, &token.UserId, &token.TokenHash, &token.ExpiresAt, &token.CreatedAt); err != nil {
		return nil, fmt.Errorf("failed to scan refresh token: %w", err)
	}
	return &token, nil
}
