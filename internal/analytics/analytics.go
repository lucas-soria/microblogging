package analytics

import "time"

// UserAnalytics represents analytics data for a user
type UserAnalytics struct {
	Handler      string    `json:"handler"`
	IsInfluencer bool      `json:"is_influencer"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Event represents an analytics event
type Event struct {
	ID        string    `json:"id"`
	EventType string    `json:"event_type"`
	Handler   string    `json:"handler"`
	TweetID   string    `json:"tweet_id,omitempty"`
	Timestamp time.Time `json:"timestamp"`
}
