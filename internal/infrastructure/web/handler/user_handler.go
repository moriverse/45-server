package handler

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	userService "github.com/moriverse/45-server/internal/app/user"
	domainErrors "github.com/moriverse/45-server/internal/domain/errors"
	"github.com/moriverse/45-server/internal/domain/user"
	"github.com/moriverse/45-server/internal/infrastructure/web/context"
	"github.com/moriverse/45-server/internal/infrastructure/web/response"
)

type UserHandler struct {
	userService userService.UserService
}

func NewUserHandler(userService userService.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.GET("/users/me", h.GetProfile)
	router.PUT("/users/me", h.UpdateProfile)
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	id, ok := context.GetUserID(c)
	if !ok {
		h.handleError(c, errors.New("failed to get user ID from context"))
		return
	}

	userProfile, err := h.userService.GetUserProfile(c.Request.Context(), id)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Data(c, http.StatusOK, userProfile)
}

type UpdateProfileRequest struct {
	Name        *string `json:"name"`
	PhoneNumber *string `json:"phoneNumber"`
	AvatarURL   *string `json:"avatarUrl"`
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	id, ok := context.GetUserID(c)
	if !ok {
		h.handleError(c, errors.New("failed to get user ID from context"))
		return
	}

	var req UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, response.APIError{
			Code:    response.CodeInvalidRequestBody,
			Message: err.Error(),
		})
		return
	}

	params := user.UpdateProfileParams{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
		AvatarURL:   req.AvatarURL,
	}

	updatedUser, err := h.userService.UpdateUserProfile(c.Request.Context(), id, params)
	if err != nil {
		h.handleError(c, err)
		return
	}

	response.Data(c, http.StatusOK, updatedUser)
}

func (h *UserHandler) handleError(c *gin.Context, err error) {
	requestLogger, ok := context.GetLogger(c)
	if !ok {
		requestLogger = slog.Default()
	}

	var notFoundErr domainErrors.NotFoundError
	if errors.As(err, &notFoundErr) {
		response.Error(c, http.StatusNotFound, response.APIError{
			Code:    response.CodeNotFound,
			Message: notFoundErr.Error(),
		})
		return
	}

	requestLogger.Error("Unhandled error in UserHandler", "error", err)
	response.Error(c, http.StatusInternalServerError, response.APIError{
		Code:    response.CodeInternalError,
		Message: "An unexpected error occurred.",
	})
}