package services

import (
	"github.com/raulito1500/merkadapp/src/bill/models"
	"github.com/raulito1500/merkadapp/src/bill/repository"
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

func (b *BillService) ListBills() []*models.Bill {
	return b.billRepository.ListBills()
}

func (b *BillService) ListBill(id string) (models.Bill, error) {
	return b.billRepository.ListBill(id)
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

func (b *BillService) UdpateBill(id string, bill *models.Bill) error {
	bill.Total = 0
	for _, i := range bill.Items {
		i.Total = i.Quantity * i.UnitValue * (1 - i.Discount)
		bill.Total += i.Total
	}
	for _, b := range bill.Bags {
		b.Total = b.Quantity * b.Value
		bill.Total += b.Total
	}
	for _, t := range bill.Taxes {
		bill.Total += t.Total
	}
	return b.billRepository.UdpateBill(id, bill)
}

func (b *BillService) MarkSpentItem(idBill string, idItem string) error {
	return b.billRepository.MarkSpentItem(idBill, idItem)
}

func (b *BillService) MergeBills(idDestination string, idsOrigen []string) error {
	return b.billRepository.MergeBills(idDestination, idsOrigen)
}
