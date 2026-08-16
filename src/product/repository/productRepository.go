package repository

import (
	"github.com/raulito1500/merkadapp/src/product/models"
)

type ProductRepository interface {
	ListProducts() []*models.ProductListItem
	InsertProduct(product *models.Product) (string, error)
	UpdateProduct(id string) error
}
