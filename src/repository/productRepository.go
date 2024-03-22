package repository

import (
	"github.com/raulito1500/merkadapp/src/models"
)

type ProductRepository interface {
	ListProducts() []*models.Product
	InsertProduct(product *models.Product) (string, error)
	UpdateProduct(id string) error
}
