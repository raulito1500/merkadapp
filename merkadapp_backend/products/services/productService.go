package services

import (
	"github.com/raulito1500/merkadapp/products/models"
	"github.com/raulito1500/merkadapp/products/repository"
)

// TODO: Validar logica de negocios
type ProductService struct {
}

func GetAll() []*models.Product {
	return repository.GetAll()
}

func UpdateProduct(id string) error {
	return repository.UpdateProduct(id)
}
