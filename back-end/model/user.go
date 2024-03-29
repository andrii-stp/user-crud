package model

// User example
type User struct {
	UserID     int64  `json:"id"`
	UserName   string `json:"user_name"   validate:"required"`
	FirstName  string `json:"first_name"  validate:"required"`
	LastName   string `json:"last_name"   validate:"required"`
	Email      string `json:"email"       validate:"required,email"`
	Status     string `json:"user_status" validate:"required,status"`
	Department string `json:"department"`
}
