package repository

import (
	"time"

	"github.com/raulito1500/merkadapp/src/bill/entities"
	"github.com/raulito1500/merkadapp/src/bill/models"
)

type BillRepository interface {
	ListBills() []*models.Bill
	ListBill(id string) (models.Bill, error)
	InsertBill(bill *models.Bill) (string, error)
	MarkSpentItem(idBill string, idItem string) error
	UdpateBill(id string, bill *models.Bill) error
	MergeBills(idDestination string, idsOrigen []string) error
	TotalByMonth(startDate time.Time, endDate time.Time) ([]*entities.BillTotal, error)
}
