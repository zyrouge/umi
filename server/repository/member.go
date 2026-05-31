package repository

type UmiMember struct {
	UserId    string        `json:"user_id"`
	TeamId    string        `json:"team_id"`
	Role      UmiMemberRole `json:"role"`
	CreatedAt int64         `json:"created_at"`
	UpdatedAt int64         `json:"updated_at"`
}
