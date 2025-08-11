package analytics

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/lucas-soria/microblogging/pkg/database"

	"gorm.io/gorm"
)

// PostgresAnalyticsRepository is a PostgreSQL implementation of the Repository interface
type PostgresAnalyticsRepository struct {
	db database.DBClient
}

// NewPostgresAnalyticsRepository creates a new PostgreSQL analytics repository
func NewPostgresAnalyticsRepository(db database.DBClient) *PostgresAnalyticsRepository {
	// Auto migrate the schema one by one
	if err := db.AutoMigrate(&UserAnalytics{}); err != nil {
		panic(fmt.Sprintf("failed to migrate UserAnalytics table: %v", err))
	}
	if err := db.AutoMigrate(&Event{}); err != nil {
		panic(fmt.Sprintf("failed to migrate Event table: %v", err))
	}

	// Create indexes if they don't exist
	if err := db.WithContext(context.Background()).Exec(`
		CREATE INDEX IF NOT EXISTS idx_events_handler ON events(handler);
		CREATE INDEX IF NOT EXISTS idx_events_timestamp ON events(timestamp);
	`).Error; err != nil {
		panic(fmt.Sprintf("failed to create database indexes: %v", err))
	}

	// Create mock user analytics
	mockUsers := []UserAnalytics{
		{
			Handler:      "lucas",
			IsInfluencer: true,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Handler:      "lucas1",
			IsInfluencer: false,
			IsActive:     true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
		{
			Handler:      "lucas2",
			IsInfluencer: false,
			IsActive:     false,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	// Save or update mock users in the database
	ctx := context.Background()
	for _, user := range mockUsers {
		// Try to find existing user
		var existing UserAnalytics
		err := db.First(ctx, &existing, "handler = ?", user.Handler)

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				// User doesn't exist, create it
				if createErr := db.Create(ctx, &user); createErr != nil {
					log.Printf("failed to create mock user %s: %v", user.Handler, createErr)
				}
			} else {
				log.Printf("error checking for mock user %s: %v", user.Handler, err)
			}
		} else {
			// User exists, update it
			existing.IsInfluencer = user.IsInfluencer
			existing.IsActive = user.IsActive
			existing.UpdatedAt = time.Now()

			if updateErr := db.Save(ctx, &existing); updateErr != nil {
				log.Printf("failed to update mock user %s: %v", user.Handler, updateErr)
			}
		}
	}

	// Print mock users from db
	var users []UserAnalytics
	if err := db.WithContext(ctx).Find(&users).Error; err != nil {
		log.Fatalf("failed to find mock users: %v", err)
	}
	for _, user := range users {
		log.Printf("mock user: %s", user.Handler)
	}

	return &PostgresAnalyticsRepository{db: db}
}

// GetUserAnalytics retrieves analytics for a specific user
func (r *PostgresAnalyticsRepository) GetUserAnalytics(ctx context.Context, userID string) (*UserAnalytics, error) {
	var analytics UserAnalytics
	if err := r.db.WithContext(ctx).Where("handler = ?", userID).First(&analytics).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user analytics not found")
		}
		return nil, fmt.Errorf("failed to get user analytics: %w", err)
	}
	return &analytics, nil
}

// GetAllUserAnalytics retrieves analytics for all users
func (r *PostgresAnalyticsRepository) GetAllUserAnalytics(ctx context.Context) ([]*UserAnalytics, error) {
	var analytics []*UserAnalytics
	if err := r.db.WithContext(ctx).Find(&analytics).Error; err != nil {
		return nil, fmt.Errorf("failed to get all user analytics: %w", err)
	}
	return analytics, nil
}

// DeleteUserAnalytics deletes analytics data for a specific user
func (r *PostgresAnalyticsRepository) DeleteUserAnalytics(ctx context.Context, userID string) error {
	if err := r.db.WithContext(ctx).Where("handler = ?", userID).Delete(&UserAnalytics{}).Error; err != nil {
		return fmt.Errorf("failed to delete user analytics: %w", err)
	}
	return nil
}

// ProcessEvent processes an analytics event
func (r *PostgresAnalyticsRepository) ProcessEvent(ctx context.Context, event *Event) error {
	// Set timestamp if not set
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}

	// Start a transaction
	tx := r.db.WithContext(ctx).Begin()
	if tx.Error != nil {
		return fmt.Errorf("failed to start transaction: %w", tx.Error)
	}

	// Save the event
	if err := tx.WithContext(ctx).Create(event).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to save event: %w", err)
	}

	// Update user analytics based on event type
	switch event.EventType {
	case "tweet_created":
		return tx.Commit().Error
	case "timeline_viewed":
		return tx.Commit().Error
	default:
		return tx.Commit().Error
	}
}
