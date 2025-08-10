package models

import "github.com/lucas-soria/microblogging/internal/tweets"

// CreateTweetRequest represents the request to create a new tweet
type CreateTweetRequest struct {
	Content Content `json:"content" binding:"required"`
}

type Content struct {
	Text string `json:"text" binding:"required"`
}

func (c *CreateTweetRequest) ToTweet() *tweets.Tweet {
	return &tweets.Tweet{
		Content: tweets.Content{
			Text: c.Content.Text,
		},
	}
}
