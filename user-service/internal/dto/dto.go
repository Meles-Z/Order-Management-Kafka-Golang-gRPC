package dto

type CreateUserDTO struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Password    string `json:"password"`
	Address     string `json:"address"`
}
type PasswordUpdateDTO struct {
	ID       string `json:"id"`
	Password string `json:"password"`
}

type UpdateUserDTO struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
