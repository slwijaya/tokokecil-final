package dto

type RegisterRequest struct {
	Name     string `json:"name" example:"Budi"`
	Email    string `json:"email" example:"budi@example.com"`
	Password string `json:"password" example:"123456"`
}
