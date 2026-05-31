package repository

type UmiTeam struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	EncryptionKey string `json:"-"`
	CreatedAt     int64  `json:"created_at"`
	UpdatedAt     int64  `json:"updated_at"`
}
