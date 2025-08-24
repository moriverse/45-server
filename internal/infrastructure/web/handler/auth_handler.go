package handler

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	authService "github.com/moriverse/45-server/internal/app/auth"
	authDomain "github.com/moriverse/45-server/internal/domain/auth"
	"github.com/moriverse/45-server/internal/infrastructure/web/middleware"
	"github.com/moriverse/45-server/internal/infrastructure/web/response"
)

// AuthHandler handles authentication-related HTTP requests.
type AuthHandler struct {
	authService *authService.Service
}

// NewAuthHandler creates a new instance of AuthHandler.
func NewAuthHandler(authService *authService.Service) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterRequest defines the flexible request body for user registration.
type RegisterRequest struct {
	Provider    string                 `json:"provider" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
	Profile     map[string]interface{} `json:"profile"` // Optional profile data
}

// Register handles the HTTP request for user registration.
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := response.APIError{
			Code:    "INVALID_REQUEST_BODY",
			Message: err.Error(),
		}
		response.Error(c, http.StatusBadRequest, apiErr)
		return
	}

	var result *authService.RegisterResult
	var err error

	switch authDomain.Provider(req.Provider) {
	case authDomain.Email:
		email, ok := req.Credentials["email"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Email is required and must be a string.",
			})
			return
		}
		password, ok := req.Credentials["password"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Password is required and must be a string.",
			})
			return
		}
		params := authService.RegisterWithEmailParams{
			Email:    email,
			Password: password,
		}
		result, err = h.authService.RegisterWithEmail(c.Request.Context(), params)

	case authDomain.Phone:
		phone, ok := req.Credentials["phone"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Phone is required and must be a string.",
			})
			return
		}
		code, ok := req.Credentials["code"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Code is required and must be a string.",
			})
			return
		}
		params := authService.RegisterWithPhoneParams{
			PhoneNumber: phone,
			Code:        code,
		}
		result, err = h.authService.RegisterWithPhone(c.Request.Context(), params)

	case authDomain.Wechat:
		code, ok := req.Credentials["code"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Code is required and must be a string.",
			})
			return
		}
		params := authService.RegisterWithWechatParams{
			Code: code,
		}
		result, err = h.authService.RegisterWithWechat(c.Request.Context(), params)

	default:
		response.Error(c, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_PROVIDER",
			Message: "The specified provider is not supported.",
		})
		return
	}

	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Data(c, http.StatusCreated, gin.H{
		"user":  result.User,
		"token": result.Token,
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	logger, _ := c.Get(middleware.LoggerKey)
	requestLogger, ok := logger.(*slog.Logger)
	if !ok {
		// Fallback to a default logger if the one in the context is not valid
		requestLogger = slog.Default()
	}

	// We check for specific, known application errors first.
	switch err {
	case authService.ErrUserAlreadyExists:
		response.Error(c, http.StatusConflict, response.APIError{
			Code:    "USER_ALREADY_EXISTS",
			Message: "A user with this identity already exists.",
		})
	default:
		// For unhandled or unexpected errors, log them and return a generic 500.
		requestLogger.Error("Unhandled API error", "error", err)
		response.Error(c, http.StatusInternalServerError, response.APIError{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "An unexpected error occurred on our end.",
		})
	}
}
