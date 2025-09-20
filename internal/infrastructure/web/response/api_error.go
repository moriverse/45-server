package response

type ErrorCode string

const (
	CodeNotFound           ErrorCode = "NOT_FOUND"
	CodeConflict           ErrorCode = "CONFLICT"
	CodeInvalidRequestBody ErrorCode = "INVALID_REQUEST_BODY"
	CodeUnauthorized       ErrorCode = "UNAUTHORIZED"
	CodeInternalError      ErrorCode = "INTERNAL_SERVER_ERROR"
	CodeNotImplemented     ErrorCode = "NOT_IMPLEMENTED"
	CodeInvalidProvider    ErrorCode = "INVALID_PROVIDER"
	CodeInvalidCredentials ErrorCode = "INVALID_CREDENTIALS"
	CodeInvalidToken       ErrorCode = "INVALID_TOKEN"
)

type APIError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}
