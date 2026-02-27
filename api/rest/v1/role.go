package restv1

// Role is the REST representation of a permission role.
type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RoleListResponse is the list of roles.
type RoleListResponse struct {
	Items []Role `json:"items"`
}
