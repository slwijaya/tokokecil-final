package dto

type ProductResponse struct {
	Name  string  `json:"name" example:"Kopi Gayo"`
	Price float64 `json:"price" example:"25000"`
}
