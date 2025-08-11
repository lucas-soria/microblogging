package tweets

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Tweet represents the tweet domain and DB model merged
type Tweet struct {
	ID        string    `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Handler   string    `gorm:"type:varchar(255);not null;index" json:"handler"`
	Content   Content   `gorm:"type:jsonb;not null" json:"content"`
	CreatedAt time.Time `gorm:"index" json:"created_at"`
}

// Content represents the content of a tweet
type Content struct {
	Text string `json:"text" validate:"max=280"`
}

// Implement driver.Valuer interface: converts Content to JSON for DB storage
func (c Content) Value() (driver.Value, error) {
	return json.Marshal(c)
}

// Implement sql.Scanner interface: converts JSON from DB to Content
func (c *Content) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan Content: type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, c)
}
