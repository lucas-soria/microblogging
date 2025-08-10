package analytics

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// assertUserAnalyticsEqual compares two slices of UserAnalytics pointers by their values regardless of order
func assertUserAnalyticsEqual(expected, actual []*UserAnalytics) func() bool {
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
				if !found[e.Handler] && e.Handler == a.Handler &&
					e.IsInfluencer == a.IsInfluencer && e.IsActive == a.IsActive {
					found[e.Handler] = true
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

func TestInMemoryRepository_GetUserAnalytics(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryRepository()

	type want struct {
		err       error
		analytics *UserAnalytics
	}

	tt := []struct {
		name         string
		expectations func()
		userID       string
		want         want
	}{
		{
			name:         "non-existent user returns error",
			expectations: func() {},
			userID:       "nonexistent-user",
			want: want{
				err:       errors.New("user analytics not found"),
				analytics: nil,
			},
		},
		{
			name: "existing user returns analytics",
			expectations: func() {
				repo.analytics["user1"] = &UserAnalytics{
					Handler:      "user1",
					IsActive:     true,
					IsInfluencer: false,
				}
			},
			userID: "user1",
			want: want{
				err: nil,
				analytics: &UserAnalytics{
					Handler:      "user1",
					IsActive:     true,
					IsInfluencer: false,
				},
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			analytics, err := repo.GetUserAnalytics(ctx, tc.userID)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.analytics, analytics)
		})
	}
}

func TestInMemoryRepository_ProcessEvent(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryRepository()

	tt := []struct {
		name  string
		event *Event
		want  error
	}{
		{
			name: "tweet_created event creates new user analytics",
			event: &Event{
				ID:        "event-1",
				EventType: "tweet_created",
				Handler:   "user1",
				Timestamp: time.Now(),
			},
			want: nil,
		},
		{
			name: "timeline_viewed event marks user as active",
			event: &Event{
				ID:        "event-2",
				EventType: "timeline_viewed",
				Handler:   "user1",
				Timestamp: time.Now(),
			},
			want: nil,
		},
		{
			name: "invalid event type is ignored",
			event: &Event{
				ID:        "event-3",
				EventType: "invalid_event",
				Handler:   "user1",
				Timestamp: time.Now(),
			},
			want: nil,
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			err := repo.ProcessEvent(ctx, tc.event)

			assert.Equal(t, tc.want, err)
		})
	}
}

func TestInMemoryRepository_GetAllUserAnalytics(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryRepository()

	type want struct {
		analytics []*UserAnalytics
		err       error
	}

	tt := []struct {
		name         string
		expectations func()
		want         want
	}{
		{
			name:         "empty repository returns empty slice",
			expectations: func() {},
			want: want{
				analytics: []*UserAnalytics{},
				err:       nil,
			},
		},
		{
			name: "returns all user analytics",
			expectations: func() {
				repo.analytics["user1"] = &UserAnalytics{Handler: "user1"}
				repo.analytics["user2"] = &UserAnalytics{Handler: "user2"}
			},
			want: want{
				analytics: []*UserAnalytics{
					{Handler: "user1"},
					{Handler: "user2"},
				},
				err: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			result, err := repo.GetAllUserAnalytics(ctx)

			assert.Condition(t, assertUserAnalyticsEqual(tc.want.analytics, result))
			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryRepository_DeleteUserAnalytics(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryRepository()

	tt := []struct {
		name         string
		expectations func()
		userID       string
		want         error
	}{
		{
			name: "deletes existing user analytics",
			expectations: func() {
				repo.analytics["user1"] = &UserAnalytics{Handler: "user1"}
			},
			userID: "user1",
			want:   nil,
		},
		{
			name:         "deleting non-existent user is not an error",
			expectations: func() {},
			userID:       "nonexistent",
			want:         errors.New("user analytics not found"),
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := repo.DeleteUserAnalytics(ctx, tc.userID)

			assert.Equal(t, tc.want, err)
		})
	}
}

func TestInMemoryRepository_UserBecomesInfluencer(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryRepository()

	userID := "user1"

	// Process exactly the threshold number of tweets (100)
	for i := 0; i < 100; i++ {
		event := &Event{
			ID:        "event-" + string(rune(i)),
			EventType: "tweet_created",
			Handler:   userID,
			Timestamp: time.Now(),
		}
		err := repo.ProcessEvent(ctx, event)
		require.NoError(t, err)
	}

	// User should not be an influencer yet
	analytics, err := repo.GetUserAnalytics(ctx, userID)
	require.NoError(t, err)
	assert.False(t, analytics.IsInfluencer)

	// Process one more tweet to cross the threshold
	event := &Event{
		ID:        "event-101",
		EventType: "tweet_created",
		Handler:   userID,
		Timestamp: time.Now(),
	}
	err = repo.ProcessEvent(ctx, event)
	require.NoError(t, err)

	// Now user should be an influencer
	analytics, err = repo.GetUserAnalytics(ctx, userID)
	require.NoError(t, err)
	assert.True(t, analytics.IsInfluencer)
}
