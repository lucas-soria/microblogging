package feed

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// assertTweetsEqual compares two slices of Tweet pointers by their values regardless of order
func assertTweetsEqual(expected, actual []*Tweet) func() bool {
	return func() bool {
		if len(expected) != len(actual) {
			return false
		}

		// Create a map to track found tweets by their IDs
		found := make(map[string]bool, len(expected))

		// For each tweet in actual, try to find a matching tweet in expected
		for _, a := range actual {
			matched := false
			for _, e := range expected {
				if !found[e.ID] && e.ID == a.ID && e.Handler == a.Handler &&
					e.Content.Text == a.Content.Text && e.CreatedAt.Equal(a.CreatedAt) {
					found[e.ID] = true
					matched = true
					break
				}
			}
			if !matched {
				return false
			}
		}

		// All tweets should be matched
		return len(found) == len(expected)
	}
}

func TestInMemoryFeedRepository_GetUserTimeline(t *testing.T) {
	ctx := context.Background()

	now := time.Now().UTC()

	type args struct {
		userID string
		limit  int
		offset int
	}

	type want struct {
		tweets []*Tweet
		err    error
	}

	tt := []struct {
		name  string
		setup func(*InMemoryFeedRepository)
		args  args
		want  want
	}{
		{
			name: "get user timeline with tweets",
			setup: func(repo *InMemoryFeedRepository) {
				t1 := &Tweet{
					ID:        "1",
					Handler:   "user1",
					Content:   Content{Text: "First tweet"},
					CreatedAt: now,
				}
				t2 := &Tweet{
					ID:        "2",
					Handler:   "user1",
					Content:   Content{Text: "Second tweet"},
					CreatedAt: now,
				}
				repo.AddTweet("user1", t1)
				repo.AddTweet("user1", t2)
			},
			args: args{
				userID: "user1",
				limit:  10,
				offset: 0,
			},
			want: want{
				tweets: []*Tweet{
					{
						ID:        "2",
						Handler:   "user1",
						Content:   Content{Text: "Second tweet"},
						CreatedAt: now,
					},
					{
						ID:        "1",
						Handler:   "user1",
						Content:   Content{Text: "First tweet"},
						CreatedAt: now,
					},
				},
				err: nil,
			},
		},
		{
			name: "pagination works correctly",
			setup: func(repo *InMemoryFeedRepository) {
				for i := 0; i < 5; i++ {
					repo.AddTweet("user1", &Tweet{
						ID:        string(rune('a' + i)),
						Handler:   "user1",
						Content:   Content{Text: string(rune('a' + i))},
						CreatedAt: now,
					})
				}
			},
			args: args{
				userID: "user1",
				limit:  2,
				offset: 1,
			},
			want: want{
				tweets: []*Tweet{
					{
						ID:        "b",
						Handler:   "user1",
						Content:   Content{Text: "b"},
						CreatedAt: now,
					},
					{
						ID:        "c",
						Handler:   "user1",
						Content:   Content{Text: "c"},
						CreatedAt: now,
					},
				},
				err: nil,
			},
		},
		{
			name:  "non-existent user returns empty",
			setup: func(*InMemoryFeedRepository) {},
			args: args{
				userID: "nonexistent",
				limit:  10,
				offset: 0,
			},
			want: want{
				tweets: []*Tweet{},
				err:    nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryFeedRepository()
			if tc.setup != nil {
				tc.setup(repo)
			}

			tweets, err := repo.GetUserTimeline(ctx, tc.args.userID, tc.args.limit, tc.args.offset)

			assert.Condition(t, assertTweetsEqual(tc.want.tweets, tweets))
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryFeedRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryFeedRepository()

	// Number of concurrent operations
	numOps := 100
	done := make(chan bool, numOps)
	errCh := make(chan error, numOps)

	// Concurrently add tweets
	for i := 0; i < numOps; i++ {
		go func(i int) {
			tweet := &Tweet{
				ID:        string(rune('a' + i)),
				Handler:   "user1",
				Content:   Content{Text: string(rune('a' + i))},
				CreatedAt: time.Now(),
			}

			repo.AddTweet("user1", tweet)

			// Also test reading concurrently
			_, err := repo.GetUserTimeline(context.Background(), "user1", 10, 0)
			errCh <- err
			done <- true
		}(i)
	}

	// Wait for all operations to complete
	for i := 0; i < numOps; i++ {
		<-done
	}

	// Check for errors
	close(errCh)
	for err := range errCh {
		assert.NoError(t, err)
	}

	// Verify final state
	timeline, err := repo.GetUserTimeline(context.Background(), "user1", numOps, 0)
	assert.NoError(t, err)
	assert.Len(t, timeline, numOps)
}
