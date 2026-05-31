package repository

type UmiEventTag struct {
	EventId   string `json:"event_id"`
	TagId     string `json:"tag_id"`
	CreatedAt int64  `json:"created_at"`
}
