package users

// User represents the user domain and DB model merged
type User struct {
	Handler   string `gorm:"primaryKey;type:varchar(255);uniqueIndex" json:"handler"`
	FirstName string `gorm:"type:varchar(255);not null" json:"first_name"`
	LastName  string `gorm:"type:varchar(255);not null" json:"last_name"`
}

// TableName specifies the table name for the User
func (User) TableName() string {
	return "users"
}

// UserFollow represents the follow relationship between users
type UserFollow struct {
	FollowerHandler string `gorm:"primaryKey;type:varchar(255);not null"`
	FolloweeHandler string `gorm:"primaryKey;type:varchar(255);not null"`
}

// TableName specifies the table name for the UserFollow
func (UserFollow) TableName() string {
	return "user_follows"
}
