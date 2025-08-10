package handlers

import (
	"net/http"

	"github.com/lucas-soria/microblogging/cmd/tweets/models"

	"github.com/lucas-soria/microblogging/internal/tweets"

	"github.com/gin-gonic/gin"
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
func (handler *TweetHandler) CreateTweet(ctx *gin.Context) {
	var tweetRequest models.CreateTweetRequest
	if err := ctx.ShouldBindJSON(&tweetRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	contentToCreate := tweetRequest.ToTweet()
	userID, _ := ctx.Get("user_id")
	contentToCreate.Handler = userID.(string)

	tweet, err := handler.service.CreateTweet(ctx.Request.Context(), contentToCreate)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tweet"})
		return
	}

	ctx.JSON(http.StatusCreated, tweet)
}

// GetTweet handles GET /v1/tweets/:id
func (handler *TweetHandler) GetTweet(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tweet ID is required"})
		return
	}

	tweet, err := handler.service.GetTweet(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tweet"})
		return
	}

	if tweet == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	ctx.JSON(http.StatusOK, tweet)
}

// GetUserTweets handles GET /v1/tweets/users/:id
func (handler *TweetHandler) GetUserTweets(ctx *gin.Context) {
	userID := ctx.Param("id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "User ID is required"})
		return
	}

	userTweets, err := handler.service.GetUserTweets(ctx.Request.Context(), userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user tweets"})
		return
	}

	if userTweets == nil {
		userTweets = []*tweets.Tweet{} // Return empty array instead of null
	}

	ctx.JSON(http.StatusOK, userTweets)
}

// DeleteTweet handles DELETE /v1/tweets/:id
func (handler *TweetHandler) DeleteTweet(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Tweet ID is required"})
		return
	}

	// Get the tweet to check ownership
	tweet, err := handler.service.GetTweet(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get tweet"})
		return
	}

	if tweet == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "Tweet not found"})
		return
	}

	// Check if the authenticated user is the owner of the tweet
	if tweet.Handler != ctx.GetHeader("X-User-Id") {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own tweets"})
		return
	}

	if err := handler.service.DeleteTweet(ctx.Request.Context(), id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete tweet"})
		return
	}

	ctx.Status(http.StatusNoContent)
}
