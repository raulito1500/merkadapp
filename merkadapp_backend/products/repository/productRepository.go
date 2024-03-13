package repository

import (
	"github.com/raulito1500/merkadapp/products/models"
)

type ProductRepository interface {
	ListProducts() []*models.Product
	UpdateProduct(id string) error
}
