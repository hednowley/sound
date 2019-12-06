package dto

// Credentials are data from a login form.
type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
