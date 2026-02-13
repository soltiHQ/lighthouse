package v1

type Credential struct {
	Auth   string `json:"auth"`
	ID     string `json:"id"`
	UserID string `json:"user_id"`
}

type CredentialListResponse struct {
	Items []Credential `json:"items"`
}
