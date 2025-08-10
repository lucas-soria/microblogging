package handlers

import (
	"net/http"

	"github.com/lucas-soria/microblogging/internal/users"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service users.UserService
}

func NewUserHandler(service users.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

type createUserRequest struct {
	Handler   string `json:"handler" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

// CreateUser handles POST /v1/users
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &users.User{
		Handler:   req.Handler,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	createdUser, err := h.service.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// GetUser handles GET /v1/users/:id
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	user, err := h.service.GetUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser handles DELETE /v1/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would check if the authenticated user has permission to delete this user
	// For now, we'll just delete the user if they exist

	if err := h.service.DeleteUser(c.Request.Context(), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// FollowUser handles POST /v1/users/:id/follow
func (h *UserHandler) FollowUser(c *gin.Context) {
	followeeID := c.Param("id")
	if followeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would get the follower ID from the authentication context
	// For now, we'll use a header
	followerID := c.GetHeader("X-User-Id")
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	req := users.FollowRequest{
		FollowerHandler: followerID,
		FolloweeHandler: followeeID,
	}

	if err := h.service.FollowUser(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}

	c.Status(http.StatusAccepted)
}

// UnfollowUser handles POST /v1/users/:id/unfollow
func (h *UserHandler) UnfollowUser(c *gin.Context) {
	followeeID := c.Param("id")
	if followeeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	// In a real app, you would get the follower ID from the authentication context
	// For now, we'll use a header
	followerID := c.GetHeader("X-User-Id")
	if followerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
		return
	}

	req := users.FollowRequest{
		FollowerHandler: followerID,
		FolloweeHandler: followeeID,
	}

	if err := h.service.UnfollowUser(c.Request.Context(), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserFollowers handles GET /v1/users/:id/followers
func (h *UserHandler) GetUserFollowers(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	followers, err := h.service.GetUserFollowers(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user followers"})
		return
	}

	c.JSON(http.StatusOK, followers)
}

// GetUserFollowing handles GET /v1/users/:id/following
func (h *UserHandler) GetUserFollowing(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user ID is required"})
		return
	}

	following, err := h.service.GetUserFollowees(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user following"})
		return
	}

	c.JSON(http.StatusOK, following)
}
