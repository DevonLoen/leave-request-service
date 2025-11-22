package entity

import (
	"time"
)

type UserRole string

const (
	RoleSuperAdmin UserRole = "superadmin"
	RoleAdmin      UserRole = "admin"
	RoleEmployee   UserRole = "employee"
)

func (r UserRole) IsValid() bool {
	switch r {
	case RoleSuperAdmin, RoleAdmin, RoleEmployee:
		return true
	}
	return false
}

type User struct {
	ID        int       `json:"id" db:"id"`
	FullName  string    `json:"fullName" db:"full_name"`
	Email     string    `json:"email" db:"email"`
	Role      UserRole  `json:"role" db:"role"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type UserFilter struct {
	Role string
}
