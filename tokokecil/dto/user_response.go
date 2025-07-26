package dto

type UserResponse struct {
	ID    uint   `json:"id" example:"1"`
	Name  string `json:"name" example:"Budi"`
	Email string `json:"email" example:"budi@example.com"`
}
