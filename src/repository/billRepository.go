package repository

import (
	"github.com/raulito1500/merkadapp/src/models"
)

type BillRepository interface {
	InsertBill(bill *models.Bill) (string, error)
}

