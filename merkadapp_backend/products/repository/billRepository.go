package repository

import (
	"github.com/raulito1500/merkadapp/products/models"
)

type BillRepository interface {
	InsertBill(bill *models.Bill) (string, error)
}

