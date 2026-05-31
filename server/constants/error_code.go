package constants

type UmiErrorCode string

const (
	UmiErrorCodeInternal     UmiErrorCode = "internal_error"
	UmiErrorCodeInvalidInput UmiErrorCode = "invalid_input"
	UmiErrorCodeNotFound     UmiErrorCode = "not_found"
	UmiErrorCodeUnauthorized UmiErrorCode = "unauthorized"
	UmiErrorCodeForbidden    UmiErrorCode = "forbidden"
	UmiErrorCodeConflict     UmiErrorCode = "conflict"
)
