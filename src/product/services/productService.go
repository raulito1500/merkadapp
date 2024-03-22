package services

import (
	"github.com/raulito1500/merkadapp/src/product/models"
	"github.com/raulito1500/merkadapp/src/product/repository"
)

type ProductService struct {
	productRepository repository.ProductRepository
}

func NewProductService(r repository.ProductRepository) ProductService {
	return ProductService{
		productRepository: r,
	}
}
func (p *ProductService) ListProducts() []*models.Product {
	return p.productRepository.ListProducts()
}

func (p *ProductService) InsertProduct(product *models.Product) (string, error) {
	product.GenMS()
	return p.productRepository.InsertProduct(product)
}

func (p *ProductService) UpdateProduct(id string) error {
	return p.productRepository.UpdateProduct(id)
}
