package repository

type UmiChannel struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	TeamId    string `json:"team_id"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
