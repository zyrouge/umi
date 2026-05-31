package repository

type UmiTag struct {
	Id        string `json:"id"`
	TeamId    string `json:"team_id"`
	Name      string `json:"name"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}
