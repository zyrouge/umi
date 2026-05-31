package repository

type UmiEventLevel string

const (
	UmiEventLevelInfo     UmiEventLevel = "info"
	UmiEventLevelWarning  UmiEventLevel = "warning"
	UmiEventLevelError    UmiEventLevel = "error"
	UmiEventLevelCritical UmiEventLevel = "critical"
)

type UmiEvent struct {
	Id        string            `json:"id"`
	ServiceId string            `json:"service_id"`
	ChannelId string            `json:"channel_id"`
	Title     string            `json:"title"`
	Body      *string           `json:"body"`
	Level     *UmiEventLevel    `json:"level" validate:"omitempty,oneof=info warning error critical"`
	ActionURL *string           `json:"action_url"`
	IconURL   *string           `json:"icon_url"`
	Metadata  map[string]string `json:"metadata"`
	CreatedAt int64             `json:"created_at"`
}
