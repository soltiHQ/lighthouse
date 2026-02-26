package restv1

// PermissionListResponse is the list of available permissions.
type PermissionListResponse struct {
	Items []string `json:"items"`
}
