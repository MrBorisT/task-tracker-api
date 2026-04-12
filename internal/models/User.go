package models

type UserRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type User struct {
	ID           string `json:"id,omitempty"`
	Email        string `json:"email,omitempty"`
	PasswordHash string `json:"-"`
}

type JWTToken struct {
	Token string `json:"token,omitempty"`
}
