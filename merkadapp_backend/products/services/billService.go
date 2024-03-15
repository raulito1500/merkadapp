package services

import (
	"github.com/raulito1500/merkadapp/products/models"
	"github.com/raulito1500/merkadapp/products/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TODO: Validar logica de negocios
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
		i.Total += i.Quantity * i.UnitValue * (1 - (i.Discount / 100))
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
