package repository

import (
	"github.com/raulito1500/merkadapp/database"
	"github.com/raulito1500/merkadapp/models"
)

// TODO: Traduce información de la BD al modelo diseñado

func GetAll() []*models.Product {
	pm := database.NewProductMongo()
	return pm.GetAll()
}

func UpdateProduct(id string) error {
	pm := database.NewProductMongo()
	return pm.UpdateProduct(id)
}
