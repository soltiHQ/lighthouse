package domain

import "time"

var _ Entity[*UserModel] = (*UserModel)(nil)

// UserModel represents a system user (human or service).
type UserModel struct {
	updatedAt time.Time
	roles     []Role
	scopes    []Scope

	id      string
	subject string
	email   string
	name    string

	disabled bool
}

// NewUserModel creates a new user domain model.
func NewUserModel(id, subject string) (*UserModel, error) {
	if id == "" {
		return nil, ErrEmptyID
	}
	if subject == "" {
		return nil, ErrInvalidSubject
	}
	return &UserModel{
		id:        id,
		subject:   subject,
		updatedAt: time.Now(),
	}, nil
}

// ID returns unique user identifier.
func (u *UserModel) ID() string {
	return u.id
}

// Subject returns stable auth subject (JWT sub).
func (u *UserModel) Subject() string {
	return u.subject
}

// Email returns user email.
func (u *UserModel) Email() string {
	return u.email
}

// Name returns display name.
func (u *UserModel) Name() string {
	return u.name
}

// Roles returns assigned roles.
func (u *UserModel) Roles() []Role {
	return append([]Role(nil), u.roles...)
}

// Scopes returns granted scopes.
func (u *UserModel) Scopes() []Scope {
	return append([]Scope(nil), u.scopes...)
}

// Disabled reports whether user is disabled.
func (u *UserModel) Disabled() bool {
	return u.disabled
}

// UpdatedAt returns last modification time.
func (u *UserModel) UpdatedAt() time.Time {
	return u.updatedAt
}

// Clone creates a deep copy.
func (u *UserModel) Clone() *UserModel {
	return &UserModel{
		id:        u.id,
		subject:   u.subject,
		email:     u.email,
		name:      u.name,
		disabled:  u.disabled,
		updatedAt: u.updatedAt,
		roles:     append([]Role(nil), u.roles...),
		scopes:    append([]Scope(nil), u.scopes...),
	}
}
