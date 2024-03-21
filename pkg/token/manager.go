package token

import (
	"time"
	"v1/pkg/model"
)

type Payload struct {
	ID       int64          `json:"id"`
	UID      string         `json:"uid"`
	Username string         `json:"username"`
	Name     string         `json:"name"`
	Role     model.RoleType `json:"role"`
}

// Manager issues token to user and verify token
type Manager interface {
	// IssueTo issues a token a User, return error if issuing process failed
	IssueTo(info Payload, expiresIn time.Duration) (string, error)

	// Verify verifies a token, and return a user info if it's a valid token, otherwise return error
	Verify(string) (Payload, error)
}
