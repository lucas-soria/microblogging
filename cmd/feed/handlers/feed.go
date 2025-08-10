package handlers

import (
	"net/http"
	"strconv"

	"github.com/lucas-soria/microblogging/internal/feed"

	"github.com/gin-gonic/gin"
)

// FeedHandler handles HTTP requests for feed operations
type FeedHandler struct {
	service feed.Service
}

// NewFeedHandler creates a new feed handler
func NewFeedHandler(service feed.Service) *FeedHandler {
	return &FeedHandler{
		service: service,
	}
}

// GetUserTimeline handles GET /v1/feed/timeline
func (handler *FeedHandler) GetUserTimeline(ctx *gin.Context) {
	// Get user ID
	user, _ := ctx.Get("user_id")
	userID := user.(string)

	// Parse query parameters
	limit, errLimit := strconv.Atoi(ctx.DefaultQuery("limit", "20"))
	if errLimit != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	offset, errOffset := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	if errOffset != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// Get user timeline
	timeline, err := handler.service.GetUserTimeline(ctx.Request.Context(), userID, limit, offset)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user timeline"})
		return
	}

	ctx.JSON(http.StatusOK, timeline)
}
