package analytics

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserAnalytics represents analytics data for a user
type UserAnalytics struct {
	Handler      string    `gorm:"primaryKey;size:64;not null" json:"handler"`
	IsInfluencer bool      `gorm:"default:false" json:"is_influencer"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `gorm:"not null;index" json:"created_at"`
	UpdatedAt    time.Time `gorm:"not null;index" json:"updated_at"`
}

// TableName specifies the table name for GORM
func (UserAnalytics) TableName() string {
	return "user_analytics"
}

// BeforeCreate is a hook that runs before creating a new record
func (u *UserAnalytics) BeforeCreate(tx *gorm.DB) error {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	if u.UpdatedAt.IsZero() {
		u.UpdatedAt = time.Now()
	}
	return nil
}

// BeforeUpdate is a hook that runs before updating a record
func (u *UserAnalytics) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// Event represents an analytics event
type Event struct {
	ID        string    `gorm:"primaryKey;size:64" json:"id"`
	EventType string    `gorm:"size:64;not null;index" json:"event_type"`
	Handler   string    `gorm:"size:64;not null;index" json:"handler"`
	TweetID   string    `gorm:"size:64;index" json:"tweet_id,omitempty"`
	Timestamp time.Time `gorm:"not null;index" json:"timestamp"`
}

// TableName specifies the table name for GORM
func (Event) TableName() string {
	return "events"
}

// BeforeCreate is a hook that runs before creating a new event
func (e *Event) BeforeCreate(tx *gorm.DB) error {
	if e.ID == "" {
		e.ID = uuid.New().String()
	}
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	return nil
}
