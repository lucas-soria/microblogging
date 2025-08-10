package tweets

import "time"

// Tweet represents a single tweet in the system
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
