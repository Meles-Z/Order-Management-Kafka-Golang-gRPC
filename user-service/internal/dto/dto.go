package dto

type PasswordUpdateRequest struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}
