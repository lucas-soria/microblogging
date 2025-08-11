package tweets

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/lucas-soria/microblogging/pkg/database"

	"gorm.io/gorm"
)

// PostgresTweetRepository is a PostgreSQL implementation of the Repository interface
type PostgresTweetRepository struct {
	db database.DBClient
}

// NewPostgresTweetRepository creates a new PostgreSQL tweet repository
func NewPostgresTweetRepository(db database.DBClient) *PostgresTweetRepository {
	// Auto migrate the schema
	if err := db.AutoMigrate(&Tweet{}); err != nil {
		log.Fatalf("failed to migrate database schema: %v", err)
	}

	// Create index on handler if it doesn't exist
	if err := db.WithContext(context.Background()).Exec(`
		CREATE INDEX IF NOT EXISTS idx_tweets_handler ON tweets(handler);
	`).Error; err != nil {
		log.Fatalf("Failed to create handler index: %v", err)
	}

	// Create mock tweets for testing
	repo := &PostgresTweetRepository{db: db}
	ctx := context.Background()

	// Sample tweet contents - each exactly 100 characters long
	tweetContents := []string{
		"Just setting up my microblogging account! So excited to be part of this amazing community of creators and thinkers! 🌟",
		"What a beautiful day to start microblogging! The sun is shining, the birds are singing, and I'm ready to share my thoughts!",
		"Hello world! My very first post here. Can't wait to connect with all the amazing people in this community. #firstpost",
		"Exploring all the fantastic features of this microblogging platform. So far, everything looks incredibly promising and exciting!",
		"Looking for some pro tips from experienced microbloggers here. What are your best practices for growing an audience?",
		"Just discovered this incredible platform and I'm already in love! The interface is so clean and user-friendly. Amazing! 👏",
		"What's everyone microblogging about today? Looking for interesting topics and discussions to join. #curious",
		"Testing, testing, 1-2-3. Is this thing on? Just making sure everything works perfectly before I start posting regularly!",
		"The community here seems absolutely fantastic! So many interesting people to follow and learn from. Truly inspiring! 💫",
		"That's all for now, folks! It's been an amazing day of microblogging. See you all tomorrow with more updates! 👋",
	}

	// Create mock tweets for each user
	users := []string{"lucas", "lucas1", "lucas2"}
	for _, user := range users {
		for i, content := range tweetContents {
			tweet := &Tweet{
				Handler: user,
				Content: Content{
					Text: fmt.Sprintf("%s - %d", content, i+1), // Add index to make tweets unique
				},
			}
			if err := repo.db.WithContext(ctx).Create(tweet).Error; err != nil {
				log.Printf("Failed to create mock tweet for %s: %v", user, err)
				continue
			}
		}
	}

	log.Printf("Created %d mock tweets for each test user", len(tweetContents))

	// Print mock tweets from db
	var tweets []*Tweet
	if err := db.WithContext(ctx).Find(&tweets).Error; err != nil {
		log.Fatalf("failed to find mock tweets: %v", err)
	}
	for _, tweet := range tweets {
		log.Printf("mock tweet: %s", tweet.Content.Text)
	}

	return repo
}

// Create saves a new tweet to the database and returns the created tweet
func (r *PostgresTweetRepository) Create(ctx context.Context, tweet *Tweet) (*Tweet, error) {
	if err := r.db.WithContext(ctx).Create(tweet).Error; err != nil {
		log.Printf("error creating tweet for handler %s: %v", tweet.Handler, err)
		return nil, err
	}

	return tweet, nil
}

// GetByID retrieves a tweet by its ID
func (r *PostgresTweetRepository) GetByID(ctx context.Context, id string) (*Tweet, error) {
	var tweet Tweet
	err := r.db.WithContext(ctx).First(&tweet, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("tweet not found with ID: %s", id)
			return nil, fmt.Errorf("tweet not found with ID: %s", id)
		}
		log.Printf("error fetching tweet with ID %s: %v", id, err)
		return nil, err
	}

	return &tweet, nil
}

// GetByUserID retrieves all tweets by a specific user
func (r *PostgresTweetRepository) GetByUserID(ctx context.Context, handler string) ([]*Tweet, error) {
	var tweets []*Tweet
	if err := r.db.WithContext(ctx).Where("handler = ?", handler).Order("created_at DESC").Find(&tweets).Error; err != nil {
		log.Printf("error fetching tweets for handler %s: %v", handler, err)
		return nil, err
	}
	return tweets, nil
}

// Delete implements the Repository interface
func (r *PostgresTweetRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&Tweet{}, "id = ?", id).Error
}
