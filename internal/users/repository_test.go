package users

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserRepository_CreateUser(t *testing.T) {

	type want struct {
		err  error
		user *User
	}

	tt := []struct {
		name  string
		user  *User
		setup func(*InMemoryUserRepository)
		want  want
	}{
		{
			name: "successful user creation",
			user: &User{
				Handler:   "testuser",
				FirstName: "Test",
				LastName:  "User",
			},
			setup: func(r *InMemoryUserRepository) {},
			want: want{
				err:  nil,
				user: &User{Handler: "testuser", FirstName: "Test", LastName: "User"},
			},
		},
		{
			name: "duplicate handler",
			user: &User{
				Handler:   "existinguser",
				FirstName: "Test",
				LastName:  "User",
			},
			setup: func(r *InMemoryUserRepository) {
				r.users["1"] = &User{Handler: "existinguser"}
			},
			want: want{
				err:  ErrHandlerExists,
				user: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			err := repo.CreateUser(context.Background(), tc.user)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				assert.NotEmpty(t, tc.user.Handler)
				assert.Equal(t, tc.want.user.Handler, tc.user.Handler)
				assert.Equal(t, tc.want.user.FirstName, tc.user.FirstName)
				assert.Equal(t, tc.want.user.LastName, tc.user.LastName)
			}
		})
	}
}

func TestInMemoryUserRepository_GetUser(t *testing.T) {

	type want struct {
		err  error
		user *User
	}

	tt := []struct {
		name    string
		handler string
		setup   func(*InMemoryUserRepository)
		want    want
	}{
		{
			name: "user found",
			setup: func(r *InMemoryUserRepository) {
				r.users[""] = &User{Handler: "testuser"}
			},
			want: want{
				err:  nil,
				user: &User{Handler: "testuser"},
			},
		},
		{
			name:    "user not found",
			handler: "nonexistent",
			setup:   func(r *InMemoryUserRepository) {},
			want: want{
				err:  ErrUserNotFound,
				user: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			user, err := repo.GetUser(context.Background(), tc.handler)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.user, user)
		})
	}
}

func TestInMemoryUserRepository_DeleteUser(t *testing.T) {

	type want struct {
		err error
	}

	tt := []struct {
		name    string
		handler string
		setup   func(*InMemoryUserRepository)
		want    want
	}{
		{
			name:    "successful deletion",
			handler: "user",
			setup: func(r *InMemoryUserRepository) {
				r.users["user"] = &User{Handler: "testuser"}
			},
			want: want{
				err: nil,
			},
		},
		{
			name:    "user not found",
			handler: "nonexistent",
			setup:   func(r *InMemoryUserRepository) {},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			err := repo.DeleteUser(context.Background(), tc.handler)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				_, err := repo.GetUser(context.Background(), tc.handler)
				assert.Equal(t, ErrUserNotFound, err)
			}
		})
	}
}

func TestInMemoryUserRepository_FollowUser(t *testing.T) {

	type want struct {
		err error
	}

	tt := []struct {
		name  string
		req   FollowRequest
		setup func(*InMemoryUserRepository)
		want  want
	}{
		{
			name: "successful follow",
			req: FollowRequest{
				FollowerHandler: "1",
				FolloweeHandler: "2",
			},
			setup: func(r *InMemoryUserRepository) {
				r.users["1"] = &User{Handler: "follower"}
				r.users["2"] = &User{Handler: "followee"}
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "follower not found",
			req: FollowRequest{
				FollowerHandler: "nonexistent",
				FolloweeHandler: "2",
			},
			setup: func(r *InMemoryUserRepository) {
				r.users["2"] = &User{Handler: "followee"}
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
		{
			name: "followee not found",
			req: FollowRequest{
				FollowerHandler: "1",
				FolloweeHandler: "nonexistent",
			},
			setup: func(r *InMemoryUserRepository) {
				r.users["1"] = &User{Handler: "follower"}
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			err := repo.FollowUser(context.Background(), tc.req)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				// Verify the follow relationship was created
				repo.mu.RLock()
				defer repo.mu.RUnlock()
				assert.True(t, repo.follow[tc.req.FollowerHandler][tc.req.FolloweeHandler])
			}
		})
	}
}

func TestInMemoryUserRepository_UnfollowUser(t *testing.T) {

	type want struct {
		err error
	}

	tt := []struct {
		name  string
		req   FollowRequest
		setup func(*InMemoryUserRepository)
		want  want
	}{
		{
			name: "successful unfollow",
			req: FollowRequest{
				FollowerHandler: "1",
				FolloweeHandler: "2",
			},
			setup: func(r *InMemoryUserRepository) {
				r.users["1"] = &User{Handler: "follower"}
				r.users["2"] = &User{Handler: "followee"}
				if r.follow["1"] == nil {
					r.follow["1"] = make(map[string]bool)
				}
				r.follow["1"]["2"] = true
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			err := repo.UnfollowUser(context.Background(), tc.req)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				// Verify the follow relationship was removed
				repo.mu.RLock()
				defer repo.mu.RUnlock()
				assert.False(t, repo.follow[tc.req.FollowerHandler][tc.req.FolloweeHandler])
			}
		})
	}
}

func TestInMemoryUserRepository_GetUserFollowers(t *testing.T) {

	type want struct {
		err   error
		users []User
	}

	tt := []struct {
		name    string
		handler string
		setup   func(*InMemoryUserRepository)
		want    want
	}{
		{
			name:    "user has followers",
			handler: "followee",
			setup: func(r *InMemoryUserRepository) {
				r.users["follower1"] = &User{Handler: "follower1"}
				r.users["followee"] = &User{Handler: "followee"}
				r.users["follower2"] = &User{Handler: "follower2"}

				r.follow["follower1"] = map[string]bool{"followee": true}
				r.follow["follower2"] = map[string]bool{"followee": true}
			},
			want: want{
				err: nil,
				users: []User{
					{Handler: "follower1"},
					{Handler: "follower2"},
				},
			},
		},
		{
			name:    "user has no followers",
			handler: "user",
			setup: func(r *InMemoryUserRepository) {
				r.users["user"] = &User{Handler: "user"}
				r.follow["user"] = map[string]bool{}
			},
			want: want{
				err:   nil,
				users: []User{},
			},
		},
		{
			name:    "user not found",
			handler: "nonexistent",
			setup:   func(r *InMemoryUserRepository) {},
			want: want{
				err:   ErrUserNotFound,
				users: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			followers, err := repo.GetUserFollowers(context.Background(), tc.handler)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				assert.ElementsMatch(t, tc.want.users, followers)
			}
		})
	}
}

func TestInMemoryUserRepository_GetUserFollowing(t *testing.T) {

	type want struct {
		err   error
		users []User
	}

	tt := []struct {
		name    string
		handler string
		setup   func(*InMemoryUserRepository)
		want    want
	}{
		{
			name:    "user is following others",
			handler: "user",
			setup: func(r *InMemoryUserRepository) {
				r.users["user"] = &User{Handler: "user"}
				r.users["followee1"] = &User{Handler: "followee1"}
				r.users["followee2"] = &User{Handler: "followee2"}

				r.follow["user"] = map[string]bool{
					"followee1": true,
					"followee2": true,
				}
			},
			want: want{
				err: nil,
				users: []User{
					{Handler: "followee1"},
					{Handler: "followee2"},
				},
			},
		},
		{
			name:    "user is not following anyone",
			handler: "user",
			setup: func(r *InMemoryUserRepository) {
				r.users["user"] = &User{Handler: "user"}
			},
			want: want{
				err:   nil,
				users: []User{},
			},
		},
		{
			name:    "user not found",
			handler: "nonexistent",
			setup:   func(r *InMemoryUserRepository) {},
			want: want{
				err:   ErrUserNotFound,
				users: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewInMemoryUserRepository()
			tc.setup(repo)

			following, err := repo.GetUserFollowing(context.Background(), tc.handler)

			assert.Equal(t, tc.want.err, err)
			if tc.want.err == nil {
				assert.ElementsMatch(t, tc.want.users, following)
			}
		})
	}
}

func TestInMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	repo := NewInMemoryUserRepository()
	user := &User{
		Handler:   "concurrent",
		FirstName: "Concurrent",
		LastName:  "User",
	}

	// Test concurrent create
	t.Run("concurrent create", func(t *testing.T) {
		for i := 0; i < 100; i++ {
			t.Run("", func(t *testing.T) {
				t.Parallel()
				u := *user
				u.Handler = "user_" + uuid.New().String()
				err := repo.CreateUser(context.Background(), &u)
				assert.NoError(t, err)
			})
		}
	})

	// Test concurrent follow
	t.Run("concurrent follow", func(t *testing.T) {
		user1 := &User{Handler: "user1"}
		user2 := &User{Handler: "user2"}
		repo.CreateUser(context.Background(), user1)
		repo.CreateUser(context.Background(), user2)

		for i := 0; i < 100; i++ {
			t.Run("", func(t *testing.T) {
				t.Parallel()
				req := FollowRequest{
					FollowerHandler: user1.Handler,
					FolloweeHandler: user2.Handler,
				}
				_ = repo.FollowUser(context.Background(), req)
				_ = repo.UnfollowUser(context.Background(), req)
			})
		}
	})
}
