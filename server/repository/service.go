package repository

type UmiService struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	TeamId    string `json:"team_id"`
	TokenHash string `json:"-"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
