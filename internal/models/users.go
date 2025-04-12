package models

import "github.com/google/uuid"

type Role string

const (
	RoleEmployee  Role = "employee"
	RoleModerator Role = "moderator"
)

type Users struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}

type DummyLogin struct {
	Role Role `json:"role"`
}

func (role Role) IsValid() bool {
	switch role {
	case RoleModerator, RoleEmployee:
		return true
	default:
		return false
	}
}
