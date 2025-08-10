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
func (h *FeedHandler) GetUserTimeline(c *gin.Context) {
	// Get user ID from header
	userID := c.GetHeader("X-User-Id")
	if userID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Id header is required"})
		return
	}

	// Parse query parameters
	limit, errLimit := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if errLimit != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit parameter"})
		return
	}
	offset, errOffset := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if errOffset != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid offset parameter"})
		return
	}

	// Get user timeline
	timeline, err := h.service.GetUserTimeline(c.Request.Context(), userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user timeline"})
		return
	}

	c.JSON(http.StatusOK, timeline)
}
