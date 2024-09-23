package users

import (
	"github.com/google/uuid"
	"go-web-api-starter/internal/validator"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	IsDeleted bool      `json:"isDeleted"`
}

// ValidateEmail checks if the provided email string is not empty and if it
// matches the regular expression for validating email addresses (EmailRX).
// This function uses the provided Validator instance for these checks.
// When the email string does not pass these checks, corresponding error
// messages are filed into the Validator's error map under 'email' key.
//
// Possible errors include:
//
//	"must be provided" - If the email is an empty string.
//	"must be a valid email address" - If the email does not match
//	the EmailRX regular expression.
func ValidateEmail(v *validator.Validator, email string) {
	v.Check(email != "", "email", "must be provided")
	v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidateUser(v *validator.Validator, user *User) {
	ValidateEmail(v, user.Email)

	if user.Role.Name == "" {
		panic("missing role for user")
	}
}
