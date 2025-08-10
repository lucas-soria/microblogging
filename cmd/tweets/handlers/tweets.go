package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lucas-soria/microblogging/internal/tweets"
)

// TweetHandler handles HTTP requests for tweet operations
type TweetHandler struct {
	service tweets.Service
}

// NewTweetHandler creates a new tweet handler
func NewTweetHandler(service tweets.Service) *TweetHandler {
	return &TweetHandler{
		service: service,
	}
}

// CreateTweet handles POST /v1/tweets
func (h *TweetHandler) CreateTweet(c *gin.Context) {
	var req tweets.CreateTweetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Get user ID from header
	req.Handler = c.GetHeader("X-User-Id")
	if req.Handler == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	tweet, err := h.service.CreateTweet(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	c.JSON(http.StatusCreated, tweet)
}

// GetTweet handles GET /v1/tweets/:id
func (h *TweetHandler) GetTweet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tweet ID is required"})
		return
	}

	// Get user ID from header
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	tweet, err := h.service.GetTweet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tweet"})
		return
	}

	if tweet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	c.JSON(http.StatusOK, tweet)
}

// GetUserTweets handles GET /v1/tweets/users/:id
func (h *TweetHandler) GetUserTweets(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	// Get user ID from header for authentication
	headerUserID := c.GetHeader("X-User-Id")
	if headerUserID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	userTweets, err := h.service.GetUserTweets(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tweets"})
		return
	}

	if userTweets == nil {
		userTweets = []*tweets.Tweet{} // Return empty array instead of null
	}

	c.JSON(http.StatusOK, userTweets)
}

// DeleteTweet handles DELETE /v1/tweets/:id
func (h *TweetHandler) DeleteTweet(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Tweet ID is required"})
		return
	}

	// Get user ID from header
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	// Get the tweet to check ownership
	tweet, err := h.service.GetTweet(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tweet"})
		return
	}

	if tweet == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	// Check if the authenticated user is the owner of the tweet
	if tweet.Handler != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own tweets"})
		return
	}

	if err := h.service.DeleteTweet(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tweet"})
		return
	}

	c.Status(http.StatusNoContent)
}
