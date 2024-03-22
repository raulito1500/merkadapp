package services

import (
	"github.com/raulito1500/merkadapp/src/models"
	"github.com/raulito1500/merkadapp/src/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BillService struct {
	billRepository repository.BillRepository
}

func NewBillService(r repository.BillRepository) BillService {
	return BillService{
		billRepository: r,
	}
}
func (b *BillService) InsertBill(bill *models.Bill) (string, error) {

	for _, i := range bill.Items {
		i.ID = primitive.NewObjectID().Hex()
		i.Total += i.Quantity * i.UnitValue * (1 - i.Discount)
		bill.Total += i.Total
	}
	for _, b := range bill.Bags {
		b.Total += b.Quantity * b.Value
		bill.Total += b.Total
	}
	for _, t := range bill.Taxes {
		bill.Total += t.Total
	}
	return b.billRepository.InsertBill(bill)
}
