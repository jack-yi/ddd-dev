package entity

import (
	"errors"
	"time"
)

const (
	StatusActive   = "active"
	StatusDisabled = "disabled"

	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	RoleOperator   = "operator"
	RoleViewer     = "viewer"
)

type User struct {
	ID        int64
	GoogleID  string
	Email     string
	Name      string
	Avatar       string
	Username     string
	PasswordHash string
	Status       string
	Roles     []Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Role struct {
	ID          int64
	Name        string
	Description string
	IsDefault   bool
	CreatedAt   time.Time
}

func NewUser(googleID, email, name, avatar string) *User {
	now := time.Now()
	return &User{
		GoogleID:  googleID,
		Email:     email,
		Name:      name,
		Avatar:    avatar,
		Status:    StatusActive,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func NewAdminUser(username, passwordHash, email string) *User {
	now := time.Now()
	return &User{
		Username:     username,
		PasswordHash: passwordHash,
		Email:        email,
		Name:         username,
		Status:       StatusActive,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

func (u *User) Disable() error {
	if u.HasRole(RoleSuperAdmin) {
		return errors.New("cannot disable super admin")
	}
	u.Status = StatusDisabled
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) Enable() {
	u.Status = StatusActive
	u.UpdatedAt = time.Now()
}

func (u *User) HasRole(roleName string) bool {
	for _, r := range u.Roles {
		if r.Name == roleName {
			return true
		}
	}
	return false
}

func (u *User) RoleNames() []string {
	names := make([]string, len(u.Roles))
	for i, r := range u.Roles {
		names[i] = r.Name
	}
	return names
}

func (u *User) IsActive() bool {
	return u.Status == StatusActive
}
