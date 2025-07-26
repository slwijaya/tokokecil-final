package repository

import (
	"tokokecil/model"

	"gorm.io/gorm"
)

type ProductRepository interface {
	GetAll() ([]model.Product, error)
	Create(product model.Product) (model.Product, error)
	Update(product model.Product) (model.Product, error)
	Delete(id uint32) error
	FindByID(id uint32) (model.Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) GetAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Find(&products).Error
	return products, err
}

func (r *productRepository) Create(product model.Product) (model.Product, error) {
	err := r.db.Create(&product).Error
	return product, err
}

func (r *productRepository) Update(product model.Product) (model.Product, error) {
	err := r.db.Save(&product).Error
	return product, err
}

func (r *productRepository) Delete(id uint32) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) FindByID(id uint32) (model.Product, error) {
	var p model.Product
	err := r.db.First(&p, id).Error
	return p, err
}
