package users

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/lucas-soria/microblogging/pkg/database"

	"gorm.io/gorm"
)

// PostgresUserRepository implements the Repository interface for PostgreSQL
type PostgresUserRepository struct {
	db database.DBClient
}

// NewPostgresUserRepository creates a new PostgreSQL user repository
func NewPostgresUserRepository(db database.DBClient) *PostgresUserRepository {
	// Auto migrate the schemas
	for _, model := range []interface{}{&User{}, &UserFollow{}} {
		if err := db.AutoMigrate(model); err != nil {
			log.Fatalf("failed to migrate database schema for %T: %v", model, err)
		}
	}

	// Create indexes if they don't exist
	if err := db.WithContext(context.Background()).Exec(`
		CREATE INDEX IF NOT EXISTS idx_users_handler ON users(handler);
		CREATE INDEX IF NOT EXISTS idx_user_follows_follower ON user_follows(follower_handler);
		CREATE INDEX IF NOT EXISTS idx_user_follows_followee ON user_follows(followee_handler);
	`).Error; err != nil {
		log.Fatalf("failed to create database indexes: %v", err)
	}
	repo := &PostgresUserRepository{db: db}
	ctx := context.Background()

	// Create mock users
	mockUsers := []*User{
		{
			Handler:   "lucas",
			FirstName: "Lucas",
			LastName:  "Soria",
		},
		{
			Handler:   "lucas1",
			FirstName: "Lucas",
			LastName:  "Soria",
		},
		{
			Handler:   "lucas2",
			FirstName: "Lucas",
			LastName:  "Soria",
		},
	}

	// Create users first
	for _, user := range mockUsers {
		// Check if user already exists
		exists, err := repo.handlerExists(ctx, user.Handler)
		if err != nil {
			log.Printf("failed to check if user %s exists: %v", user.Handler, err)
			continue
		}

		if exists {
			log.Printf("user %s already exists, skipping creation", user.Handler)
			continue
		}

		if err := repo.db.WithContext(ctx).Create(user).Error; err != nil {
			log.Printf("failed to create mock user %s: %v", user.Handler, err)
			continue
		}
		log.Printf("created mock user: %s", user.Handler)
	}

	// Set up follower relationships
	relationships := []struct {
		follower string
		followee string
	}{
		{"lucas1", "lucas"},   // lucas1 follows lucas
		{"lucas2", "lucas"},   // lucas2 follows lucas
		{"lucas", "lucas1"},   // lucas follows lucas1
		{"lucas1", "lucas2"},  // lucas1 follows lucas2
	}

	for _, rel := range relationships {
		// Check if relationship already exists
		var count int64
		err := repo.db.WithContext(ctx).Model(&UserFollow{}).
			Where("follower_handler = ? AND followee_handler = ?", rel.follower, rel.followee).
			Count(&count).Error

		if err != nil {
			log.Printf("failed to check follow relationship %s -> %s: %v", rel.follower, rel.followee, err)
			continue
		}

		if count > 0 {
			log.Printf("follow relationship %s -> %s already exists", rel.follower, rel.followee)
			continue
		}

		// Create the follow relationship
		follow := UserFollow{
			FollowerHandler: rel.follower,
			FolloweeHandler: rel.followee,
		}

		if err := repo.db.WithContext(ctx).Create(&follow).Error; err != nil {
			log.Printf("failed to create follow relationship %s -> %s: %v", rel.follower, rel.followee, err)
			continue
		}
		log.Printf("created follow relationship: %s -> %s", rel.follower, rel.followee)
	}

	// Print all users from db for verification
	var dbUsers []User
	if err := db.WithContext(ctx).Find(&dbUsers).Error; err != nil {
		log.Printf("failed to list users: %v", err)
	} else {
		log.Printf("total users in database: %d", len(dbUsers))
		for _, u := range dbUsers {
			log.Printf("user: %s (%s %s)", u.Handler, u.FirstName, u.LastName)
		}
	}

	return repo
}

// CreateUser implements the Repository interface
func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *User) error {
	// Check if handler already exists
	exists, err := r.handlerExists(ctx, user.Handler)
	if err != nil {
		log.Printf("error checking if handler %s exists: %v", user.Handler, err)
		return err
	}
	if exists {
		log.Printf("attempted to create user with existing handler: %s", user.Handler)
		return ErrHandlerExists
	}

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("error creating user with handler %s: %v", user.Handler, err)
		return err
	}

	return nil
}

// GetUser implements the Repository interface
func (r *PostgresUserRepository) GetUser(ctx context.Context, handler string) (*User, error) {
	var user User

	err := r.db.WithContext(ctx).First(&user, "handler = ?", handler).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("user not found with handler: %s", handler)
			return nil, ErrUserNotFound
		}
		log.Printf("error fetching user with handler %s: %v", handler, err)
		return nil, err
	}

	return &user, nil
}

// DeleteUser implements the Repository interface
func (r *PostgresUserRepository) DeleteUser(ctx context.Context, handler string) error {
	result := r.db.WithContext(ctx).Where("handler = ?", handler).Delete(&User{})
	if result.Error != nil {
		log.Printf("error deleting user with handler %s: %v", handler, result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		log.Printf("attempted to delete non-existent user with handler: %s", handler)
		return ErrUserNotFound
	}

	return nil
}

// FollowUser implements the Repository interface
func (r *PostgresUserRepository) FollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	// Check if both users exist
	if _, err := r.GetUser(ctx, followerHandler); err != nil {
		log.Printf("error verifying follower %s: %v", followerHandler, err)
		return fmt.Errorf("failed to verify follower: %w", err)
	}

	if _, err := r.GetUser(ctx, followeeHandler); err != nil {
		log.Printf("error verifying followee %s: %v", followeeHandler, err)
		return fmt.Errorf("failed to verify followee: %w", err)
	}

	// Create follow relationship
	follow := UserFollow{
		FollowerHandler: followerHandler,
		FolloweeHandler: followeeHandler,
	}

	if err := r.db.WithContext(ctx).Create(&follow).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			log.Printf("user %s attempted to follow %s again", followerHandler, followeeHandler)
			return fmt.Errorf("already following this user")
		}
		log.Printf("error creating follow relationship %s -> %s: %v", followerHandler, followeeHandler, err)
		return fmt.Errorf("failed to create follow relationship: %w", err)
	}

	return nil
}

// UnfollowUser implements the Repository interface
func (r *PostgresUserRepository) UnfollowUser(ctx context.Context, followerHandler string, followeeHandler string) error {
	result := r.db.WithContext(ctx).
		Where("follower_handler = ? AND followee_handler = ?", followerHandler, followeeHandler).
		Delete(&UserFollow{})

	if result.Error != nil {
		log.Printf("error unfollowing %s -> %s: %v", followerHandler, followeeHandler, result.Error)
		return fmt.Errorf("failed to unfollow user: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		log.Printf("no active follow relationship found: %s -> %s", followerHandler, followeeHandler)
		return fmt.Errorf("not following this user")
	}

	return nil
}

// GetUserFollowers implements the Repository interface
func (r *PostgresUserRepository) GetUserFollowers(ctx context.Context, followeeHandler string) ([]User, error) {
	var followers []User

	err := r.db.WithContext(ctx).Raw(`
		SELECT u.*
		FROM users u
		JOIN user_follows uf ON u.handler = uf.follower_handler
		WHERE uf.followee_handler = ?
	`, followeeHandler).Scan(&followers).Error

	if err != nil {
		return nil, err
	}

	return followers, nil
}

// GetUserFollowees implements the Repository interface
func (r *PostgresUserRepository) GetUserFollowees(ctx context.Context, followerHandler string) ([]User, error) {
	var followees []User

	err := r.db.WithContext(ctx).Raw(`
		SELECT u.*
		FROM users u
		JOIN user_follows uf ON u.handler = uf.followee_handler
		WHERE uf.follower_handler = ?
	`, followerHandler).Scan(&followees).Error

	if err != nil {
		return nil, err
	}

	return followees, nil
}

// handlerExists checks if a user with the given handler exists
func (r *PostgresUserRepository) handlerExists(ctx context.Context, handler string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&User{}).Where("handler = ?", handler).Count(&count).Error
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
