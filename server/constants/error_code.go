package constants

type UmiErrorCode string

const (
	ErrorCodeInternal     UmiErrorCode = "internal_error"
	ErrorCodeInvalidInput UmiErrorCode = "invalid_input"
	ErrorCodeNotFound     UmiErrorCode = "not_found"
	ErrorCodeUnauthorized UmiErrorCode = "unauthorized"
	ErrorCodeForbidden    UmiErrorCode = "forbidden"
	ErrorCodeConflict     UmiErrorCode = "conflict"
)
