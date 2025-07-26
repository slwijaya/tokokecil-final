package dto

type LoginRequest struct {
	Email    string `json:"email" example:"budi@example.com"`
	Password string `json:"password" example:"123456"`
}
