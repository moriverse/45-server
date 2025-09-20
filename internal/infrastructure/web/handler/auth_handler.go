package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	authService "github.com/moriverse/45-server/internal/app/auth"
	authDomain "github.com/moriverse/45-server/internal/domain/auth"
	domainErrors "github.com/moriverse/45-server/internal/domain/errors"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/web/context"
	"github.com/moriverse/45-server/internal/infrastructure/web/response"
)

type AuthHandler struct {
	authService authService.AuthService
}

func NewAuthHandler(authService authService.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/login", h.Login)
}

type LoginRequest struct {
	Provider    string                 `json:"provider" binding:"required"`
	Credentials map[string]interface{} `json:"credentials" binding:"required"`
	Source      user.Source            `json:"source"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.APIError{
			Code:    response.CodeInvalidRequestBody,
			Message: err.Error(),
		})
		return
	}

	provider := authDomain.Provider(req.Provider)

	var result *authService.LoginOrRegisterResult
	var err error

	switch provider {
	case authDomain.Wechat:
		code, ok := req.Credentials["code"].(string)
		if !ok {
			response.Error(c, http.StatusBadRequest, response.APIError{
				Code:    response.CodeInvalidCredentials,
				Message: "Wechat provider requires a 'code' string credential.",
			})
			return
		}
		params := authService.LoginOrRegisterWithWechatParams{
			Code:   code,
			Source: req.Source,
		}
		result, err = h.authService.LoginOrRegisterWithWechat(c.Request.Context(), params)

	default:
		response.Error(c, http.StatusBadRequest, response.APIError{
			Code:    response.CodeInvalidProvider,
			Message: "The specified provider is not supported.",
		})
		return
	}

	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Data(c, http.StatusOK, gin.H{
		"user":    result.User,
		"token":   result.Token,
		"newUser": result.NewUser,
	})
}

func (h *AuthHandler) handleError(c *gin.Context, err error) {
	requestLogger, ok := context.GetLogger(c)
	if !ok {
		requestLogger = slog.Default()
	}

	var dataInconsistencyErr domainErrors.DataInconsistencyError
	if errors.As(err, &dataInconsistencyErr) {
		requestLogger.Error("Critical data inconsistency detected", "error", err)
		response.Error(c, http.StatusInternalServerError, response.APIError{
			Code:    response.CodeInternalError,
			Message: "An unexpected error occurred during login.",
		})
		return
	}

	requestLogger.Warn("Login failed", "error", err)
	response.Error(c, http.StatusInternalServerError, response.APIError{
		Code:    response.CodeInternalError,
		Message: "An unexpected error occurred during login.",
	})
}
