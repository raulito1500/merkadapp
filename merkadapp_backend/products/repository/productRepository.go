package repository

import (
	"github.com/raulito1500/merkadapp/products/models"
)

// TODO: Traduce información de la BD al modelo diseñado

func GetAll() []*models.Product {
	pm := NewProductMongo()
	return pm.ListProducts()
}

func UpdateProduct(id string) error {
	pm := NewProductMongo()
	return pm.UpdateProduct(id)
}
