package service

import (
	"errors"
	"tokokecil/model"
	"tokokecil/repository"
)

type ProductService interface {
	GetAllProducts() ([]model.Product, error)
	CreateProduct(product model.Product) (model.Product, error)
	UpdateProduct(product model.Product) (model.Product, error)
	DeleteProduct(id uint32) error
	FindProductByID(id uint32) (model.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) ProductService {
	return &productService{repo: r}
}

func (s *productService) GetAllProducts() ([]model.Product, error) {
	return s.repo.GetAll()
}

func (s *productService) CreateProduct(product model.Product) (model.Product, error) {
	if product.Name == "" || product.Price <= 0 {
		return model.Product{}, errors.New("name & price wajib diisi")
	}
	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(product model.Product) (model.Product, error) {
	_, err := s.repo.FindByID(uint32(product.ID))
	if err != nil {
		return model.Product{}, errors.New("product not found")
	}
	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id uint32) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}
	return s.repo.Delete(id)
}

func (s *productService) FindProductByID(id uint32) (model.Product, error) {
	return s.repo.FindByID(id)
}
