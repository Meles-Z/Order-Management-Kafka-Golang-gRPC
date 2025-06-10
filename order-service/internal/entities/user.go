package entities

type User struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	IsActive    bool   `json:"isActive"`
}
