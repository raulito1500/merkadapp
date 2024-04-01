package repository

import (
	"github.com/raulito1500/merkadapp/src/bill/models"
)

type BillRepository interface {
	InsertBill(bill *models.Bill) (string, error)
	MarkSpentItem(idBill string, idItem string) error
}

