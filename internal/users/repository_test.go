package users

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryUserRepository_CreateUser(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err  error
		user *User
	}

	tt := []struct {
		name         string
		user         *User
		expectations func()
		want         want
	}{
		{
			name: "successful user creation",
			user: &User{
				Handler:   "testuser",
				FirstName: "Test",
				LastName:  "User",
			},
			expectations: func() {},
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
			expectations: func() {
				repo.users["existinguser"] = &User{Handler: "existinguser"}
			},
			want: want{
				err:  ErrHandlerExists,
				user: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := repo.CreateUser(ctx, tc.user)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryUserRepository_GetUser(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err  error
		user *User
	}

	tt := []struct {
		name         string
		handler      string
		expectations func()
		want         want
	}{
		{
			name: "user found",
			expectations: func() {
				repo.users[""] = &User{Handler: "testuser"}
			},
			want: want{
				err:  nil,
				user: &User{Handler: "testuser"},
			},
		},
		{
			name:         "user not found",
			handler:      "nonexistent",
			expectations: func() {},
			want: want{
				err:  ErrUserNotFound,
				user: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			user, err := repo.GetUser(ctx, tc.handler)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.user, user)
		})
	}
}

func TestInMemoryUserRepository_DeleteUser(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		handler      string
		expectations func()
		want         want
	}{
		{
			name:    "successful deletion",
			handler: "user",
			expectations: func() {
				repo.users["user"] = &User{Handler: "testuser"}
			},
			want: want{
				err: nil,
			},
		},
		{
			name:         "user not found",
			handler:      "nonexistent",
			expectations: func() {},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := repo.DeleteUser(ctx, tc.handler)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryUserRepository_FollowUser(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		follower     string
		followee     string
		expectations func()
		want         want
	}{
		{
			name:     "successful follow",
			follower: "1",
			followee: "2",
			expectations: func() {
				repo.users["1"] = &User{Handler: "follower"}
				repo.users["2"] = &User{Handler: "followee"}
			},
			want: want{
				err: nil,
			},
		},
		{
			name:     "follower not found",
			follower: "nonexistent",
			followee: "2",
			expectations: func() {
				repo.users["2"] = &User{Handler: "followee"}
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
		{
			name:     "followee not found",
			follower: "1",
			followee: "nonexistent",
			expectations: func() {
				repo.users["1"] = &User{Handler: "follower"}
			},
			want: want{
				err: ErrUserNotFound,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := repo.FollowUser(ctx, tc.follower, tc.followee)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryUserRepository_UnfollowUser(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err error
	}

	tt := []struct {
		name         string
		follower     string
		followee     string
		expectations func()
		want         want
	}{
		{
			name:     "successful unfollow",
			follower: "1",
			followee: "2",
			expectations: func() {
				repo.users["1"] = &User{Handler: "follower"}
				repo.users["2"] = &User{Handler: "followee"}
				if repo.follow["1"] == nil {
					repo.follow["1"] = make(map[string]bool)
				}
				repo.follow["1"]["2"] = true
			},
			want: want{
				err: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			err := repo.UnfollowUser(ctx, tc.follower, tc.followee)

			assert.Equal(t, tc.want.err, err)
		})
	}
}

func TestInMemoryUserRepository_GetUserFollowers(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err   error
		users []User
	}

	tt := []struct {
		name         string
		handler      string
		expectations func()
		want         want
	}{
		{
			name:    "user has followers",
			handler: "followee",
			expectations: func() {
				repo.users["follower1"] = &User{Handler: "follower1"}
				repo.users["followee"] = &User{Handler: "followee"}
				repo.users["follower2"] = &User{Handler: "follower2"}

				repo.follow["follower1"] = map[string]bool{"followee": true}
				repo.follow["follower2"] = map[string]bool{"followee": true}
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
			expectations: func() {
				repo.users["user"] = &User{Handler: "user"}
				repo.follow["user"] = map[string]bool{}
			},
			want: want{
				err:   nil,
				users: nil,
			},
		},
		{
			name:         "user not found",
			handler:      "nonexistent",
			expectations: func() {},
			want: want{
				err:   ErrUserNotFound,
				users: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			followers, err := repo.GetUserFollowers(ctx, tc.handler)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.users, followers)
		})
	}
}

func TestInMemoryUserRepository_GetUserFollowing(t *testing.T) {
	ctx := context.Background()

	repo := NewInMemoryUserRepository()

	type want struct {
		err   error
		users []User
	}

	tt := []struct {
		name         string
		handler      string
		expectations func()
		want         want
	}{
		{
			name:    "user is following others",
			handler: "user",
			expectations: func() {
				repo = NewInMemoryUserRepository()
				repo.users["user"] = &User{Handler: "user"}
				repo.users["followee1"] = &User{Handler: "followee1"}
				repo.users["followee2"] = &User{Handler: "followee2"}

				repo.follow["user"] = map[string]bool{
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
			expectations: func() {
				repo = NewInMemoryUserRepository()
				repo.users["user"] = &User{Handler: "user"}
			},
			want: want{
				err:   nil,
				users: nil,
			},
		},
		{
			name:         "user not found",
			handler:      "nonexistent",
			expectations: func() {},
			want: want{
				err:   ErrUserNotFound,
				users: nil,
			},
		},
	}
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tc.expectations()

			following, err := repo.GetUserFollowees(ctx, tc.handler)

			assert.Equal(t, tc.want.err, err)
			assert.Equal(t, tc.want.users, following)
		})
	}
}

func TestInMemoryUserRepository_ConcurrentAccess(t *testing.T) {
	ctx := context.Background()

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
				err := repo.CreateUser(ctx, &u)
				assert.NoError(t, err)
			})
		}
	})

	// Test concurrent follow
	t.Run("concurrent follow", func(t *testing.T) {
		user1 := &User{Handler: "user1"}
		user2 := &User{Handler: "user2"}
		repo.CreateUser(ctx, user1)
		repo.CreateUser(ctx, user2)

		for i := 0; i < 100; i++ {
			t.Run("", func(t *testing.T) {
				t.Parallel()
				_ = repo.FollowUser(context.Background(), user1.Handler, user2.Handler)
				_ = repo.UnfollowUser(context.Background(), user1.Handler, user2.Handler)
			})
		}
	})
}
