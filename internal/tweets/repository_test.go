package tweets

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryTweetRepository_Create(t *testing.T) {
	type want struct {
		err   error
		tweet *Tweet
	}

	tt := []struct {
		name  string
		setup func(*InMemoryTweetRepository)
		tweet *Tweet
		want  want
	}{
		{
			name:  "successful tweet creation",
			setup: func(r *InMemoryTweetRepository) {},
			tweet: &Tweet{
				ID:        uuid.NewString(),
				Handler:   "testuser",
				Content:   Content{Text: "Hello, world!"},
				CreatedAt: time.Now().UTC(),
			},
			want: want{
				err:   nil,
				tweet: &Tweet{Handler: "testuser", Content: Content{Text: "Hello, world!"}},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryTweetRepository()
			tc.setup(repo)

			created, err := repo.Create(context.Background(), tc.tweet)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				assert.NotEmpty(t, created.ID)
				assert.Equal(t, tc.want.tweet.Handler, created.Handler)
				assert.Equal(t, tc.want.tweet.Content, created.Content)
				assert.False(t, created.CreatedAt.IsZero())
			}
		})
	}
}

func TestInMemoryTweetRepository_GetByID(t *testing.T) {
	type want struct {
		err   error
		tweet *Tweet
	}

	tt := []struct {
		name  string
		setup func(*InMemoryTweetRepository)
		id    string
		want  want
	}{
		{
			name: "tweet found",
			setup: func(r *InMemoryTweetRepository) {
				tweet := &Tweet{
					ID:        "123",
					Handler:   "testuser",
					Content:   Content{Text: "Hello"},
					CreatedAt: time.Now().UTC(),
				}
				r.tweets["123"] = tweet
			},
			id: "123",
			want: want{
				err:   nil,
				tweet: &Tweet{ID: "123", Handler: "testuser", Content: Content{Text: "Hello"}},
			},
		},
		{
			name:  "tweet not found",
			setup: func(r *InMemoryTweetRepository) {},
			id:    "nonexistent",
			want: want{
				err:   nil,
				tweet: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryTweetRepository()
			tc.setup(repo)

			tweet, err := repo.GetByID(context.Background(), tc.id)

			assert.Equal(t, tc.want.err, err)
			if tc.want.tweet != nil {
				assert.Equal(t, tc.want.tweet.ID, tweet.ID)
				assert.Equal(t, tc.want.tweet.Handler, tweet.Handler)
				assert.Equal(t, tc.want.tweet.Content, tweet.Content)
			} else {
				assert.Nil(t, tweet)
			}
		})
	}
}

func TestInMemoryTweetRepository_GetByUserID(t *testing.T) {
	type want struct {
		err    error
		tweets []*Tweet
	}

	tt := []struct {
		name   string
		setup  func(*InMemoryTweetRepository)
		userID string
		want   want
	}{
		{
			name: "tweets found",
			setup: func(r *InMemoryTweetRepository) {
				r.tweets["1"] = &Tweet{ID: "1", Handler: "user1", Content: Content{Text: "Tweet 1"}}
				r.tweets["2"] = &Tweet{ID: "2", Handler: "user1", Content: Content{Text: "Tweet 2"}}
				r.tweets["3"] = &Tweet{ID: "3", Handler: "user2", Content: Content{Text: "Tweet 3"}}
			},
			userID: "user1",
			want: want{
				err: nil,
				tweets: []*Tweet{
					{ID: "1", Handler: "user1", Content: Content{Text: "Tweet 1"}},
					{ID: "2", Handler: "user1", Content: Content{Text: "Tweet 2"}},
				},
			},
		},
		{
			name:   "no tweets found",
			setup:  func(r *InMemoryTweetRepository) {},
			userID: "nonexistent",
			want: want{
				err:    nil,
				tweets: []*Tweet{},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryTweetRepository()
			tc.setup(repo)

			tweets, err := repo.GetByUserID(context.Background(), tc.userID)

			assert.Equal(t, tc.want.err, err)
			assert.Len(t, tweets, len(tc.want.tweets))
			for i, wantTweet := range tc.want.tweets {
				assert.Equal(t, wantTweet.ID, tweets[i].ID)
				assert.Equal(t, wantTweet.Handler, tweets[i].Handler)
				assert.Equal(t, wantTweet.Content, tweets[i].Content)
			}
		})
	}
}

func TestInMemoryTweetRepository_Delete(t *testing.T) {
	type want struct {
		err error
	}

	tt := []struct {
		name  string
		setup func(*InMemoryTweetRepository)
		id    string
		want  want
	}{
		{
			name: "successful deletion",
			setup: func(r *InMemoryTweetRepository) {
				r.tweets["123"] = &Tweet{ID: "123", Handler: "testuser"}
			},
			id: "123",
			want: want{
				err: nil,
			},
		},
		{
			name:  "tweet not found",
			setup: func(r *InMemoryTweetRepository) {},
			id:    "nonexistent",
			want: want{
				err: nil, // Deleting a non-existent tweet doesn't return an error
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryTweetRepository()
			tc.setup(repo)

			err := repo.Delete(context.Background(), tc.id)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil && tc.id != "nonexistent" {
				_, err := repo.GetByID(context.Background(), tc.id)
				assert.Nil(t, err) // Should be nil because the tweet was deleted
			}
		})
	}
}

func TestInMemoryTweetRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryTweetRepository()
	ctx := context.Background()

	// Number of concurrent operations
	count := 100
	done := make(chan bool, count)

	// Concurrent creates
	for i := 0; i < count; i++ {
		go func(i int) {
			tweet := &Tweet{
				ID:        uuid.NewString(),
				Handler:   "user1",
				Content:   Content{Text: "Concurrent test"},
				CreatedAt: time.Now().UTC(),
			}
			_, err := repo.Create(ctx, tweet)
			assert.NoError(t, err)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < count; i++ {
		<-done
	}

	// Verify all tweets were created
	tweets, err := repo.GetByUserID(ctx, "user1")
	assert.NoError(t, err)
	assert.Len(t, tweets, count)
}
