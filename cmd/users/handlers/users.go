package handlers

import (
	"net/http"

	"github.com/lucas-soria/microblogging/internal/users"

	"github.com/lucas-soria/microblogging/cmd/users/models"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service users.Service
}

func NewUserHandler(service users.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// CreateUser handles POST /v1/users
func (handler *UserHandler) CreateUser(ctx *gin.Context) {
	var userToCreate models.CreateUserRequest
	if err := ctx.ShouldBindJSON(&userToCreate); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	user := userToCreate.ToUser()

	createdUser, err := handler.service.CreateUser(ctx.Request.Context(), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	ctx.JSON(http.StatusCreated, createdUser)
}

// GetUser handles GET /v1/users/:id
func (handler *UserHandler) GetUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := handler.service.GetUser(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /v1/users/:id
func (handler *UserHandler) DeleteUser(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would check if the authenticated user has permission to delete this user
	// For now, we'll just delete the user if they exist

	if err := handler.service.DeleteUser(ctx.Request.Context(), userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// FollowUser handles POST /v1/users/:id/follow
func (handler *UserHandler) FollowUser(ctx *gin.Context) {
	followeeID := ctx.Param("id")
	if followeeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would get the follower ID from the authentication context
	// For now, we'll use a header
	followerID, _ := ctx.Get("user_id")
	followerIDString := followerID.(string)

	if err := handler.service.FollowUser(ctx.Request.Context(), followerIDString, followeeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}

	ctx.Status(http.StatusAccepted)
}

// UnfollowUser handles POST /v1/users/:id/unfollow
func (handler *UserHandler) UnfollowUser(ctx *gin.Context) {
	followeeID := ctx.Param("id")
	if followeeID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would get the follower ID from the authentication context
	// For now, we'll use a header
	followerID := ctx.GetHeader("X-User-Id")
	if followerID == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	if err := handler.service.UnfollowUser(ctx.Request.Context(), followerID, followeeID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

// GetUserFollowers handles GET /v1/users/:id/followers
func (handler *UserHandler) GetUserFollowers(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	followers, err := handler.service.GetUserFollowers(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user followers"})
		return
	}

	ctx.JSON(http.StatusOK, followers)
}

// GetUserFollowees handles GET /v1/users/:id/followees
func (handler *UserHandler) GetUserFollowees(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	following, err := handler.service.GetUserFollowees(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user following"})
		return
	}

	ctx.JSON(http.StatusOK, following)
}
