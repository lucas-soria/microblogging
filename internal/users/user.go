package users

type User struct {
	Handler   string `json:"handler"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
