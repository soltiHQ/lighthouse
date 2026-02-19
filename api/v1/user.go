package v1

type User struct {
	Permissions []string `json:"permissions,omitempty"`
	RoleIDs     []string `json:"role_ids,omitempty"`

	Subject string `json:"subject"`
	Email   string `json:"email,omitempty"`
	Name    string `json:"name,omitempty"`
	ID      string `json:"id"`

	Disabled bool `json:"disabled"`
}

type UserListResponse struct {
	Items      []User `json:"items"`
	NextCursor string `json:"next_cursor,omitempty"`
}
