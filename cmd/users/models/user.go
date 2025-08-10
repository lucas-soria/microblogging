package models

import "github.com/lucas-soria/microblogging/internal/users"

type CreateUserRequest struct {
	Handler   string `json:"handler" binding:"required"`
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

func (c *CreateUserRequest) ToUser() *users.User {
	return &users.User{
		Handler:   c.Handler,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
}
