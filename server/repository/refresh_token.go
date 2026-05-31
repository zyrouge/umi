package repository

import "time"

type UmiRefreshToken struct {
	Id        string `json:"id"`
	UserId    string `json:"user_id"`
	TokenHash string `json:"-"`
	ExpiresAt int64  `json:"expires_at"`
	CreatedAt int64  `json:"created_at"`
}

func (token *UmiRefreshToken) IsValid() bool {
	nowUnix := time.Now().Unix()
	return token.ExpiresAt > nowUnix
}
