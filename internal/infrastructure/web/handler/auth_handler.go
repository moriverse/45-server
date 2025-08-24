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

// LoginRequest defines the flexible request body for user login.
type LoginRequest struct {
	Provider    string                 `json:"provider" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
}

// Login handles the HTTP request for user login or seamless registration.
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.APIError{
			Code:    "INVALID_REQUEST_BODY",
			Message: err.Error(),
		})
		return
	}

	provider := authDomain.Provider(req.Provider)

	var result *authService.RegisterResult // Login and Register return the same result
	var err error

	switch provider {
	case authDomain.Wechat:
		code, ok := req.Credentials["code"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    "INVALID_CREDENTIALS",
				Message: "Code is required and must be a string.",
			})
			return
		}
		params := authService.LoginOrRegisterWithWechatParams{
			Code: code,
			// Source can be set here
		}
		result, err = h.authService.LoginOrRegisterWithWechat(c.Request.Context(), params)

	// TODO: Add case for phone login
	case authDomain.Phone:
		response.Error(c, http.StatusNotImplemented, response.APIError{
			Code:    "NOT_IMPLEMENTED",
			Message: "This login provider is not yet implemented.",
		})
		return

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

	response.Data(c, http.StatusOK, gin.H{
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