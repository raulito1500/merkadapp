package services

import (
	"github.com/raulito1500/merkadapp/products/models"
	"github.com/raulito1500/merkadapp/products/repository"
)

// TODO: Validar logica de negocios
type BillService struct {
	billRepository repository.BillRepository
}

func NewBillService(r repository.BillRepository) BillService{
	return BillService{
		billRepository: r,
	}
}
func (b *BillService) Insert(bill *models.Bill) (string, error) {
	return b.billRepository.InsertBill(bill)
}