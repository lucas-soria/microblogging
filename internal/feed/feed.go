package feed

import (
	"time"
)

// Tweet represents a tweet in the user's feed
type Tweet struct {
	ID        string    `json:"id"`
	Handler   string    `json:"handler"`
	Content   Content   `json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// Content represents the content of a tweet
type Content struct {
	Text string `json:"text"`
}

// TimelineResponse represents the response for the user timeline
type TimelineResponse struct {
	Tweets     []*Tweet `json:"tweets"`
	NextOffset int      `json:"next_offset"`
}
