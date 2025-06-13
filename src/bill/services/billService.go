package services

import (
	"encoding/xml"
	"fmt"
	"mime/multipart"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/raulito1500/merkadapp/src/bill/entities"
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

func (b *BillService) TotalByMonth() ([]*entities.BillTotal, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month()-5, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 6, 0)

	return b.billRepository.TotalByMonth(startDate, endDate)
}

func (b *BillService) BillItemsByProduct(id string) ([]*entities.BillItem, error) {
	return b.billRepository.BillItemsByProduct(id)
}

func (b *BillService) RecommendedProducts() ([]*entities.Recommendation, error) {
	return b.billRepository.RecommendedProducts()
}

func (b *BillService) UploadXML(file *multipart.FileHeader) (string, error) {
	fileContent, err := file.Open()
	if err != nil {
		return "", err
	}
	defer fileContent.Close()

	buf := make([]byte, file.Size)
	_, err = fileContent.Read(buf)
	if err != nil {
		return "", err
	}

	var attachedDocument entities.AttachedDocument
	err = xml.Unmarshal(buf, &attachedDocument)
	if err != nil {
		return "", err
	}

	invoiceXML := strings.TrimSpace(attachedDocument.Description)
	if invoiceXML == "" {
		return "", fmt.Errorf("No se encontró contenido en el XML")
	}

	invoiceXML = extractCDATA(invoiceXML)
	if !strings.HasPrefix(invoiceXML, "<?xml") {
		return "", fmt.Errorf("El contenido extraído no parece ser un XML válido.")
	}

	// Parsear el XML de la factura
	var invoice entities.Invoice
	err = xml.Unmarshal([]byte(invoiceXML), &invoice)
	if err != nil {
		return "", err
	}

	bill := new(models.Bill)
	loc, err := time.LoadLocation("America/Bogota")
	if err != nil {
		return "", fmt.Errorf("failed to load location: %v", err)
	}
	date, err := time.ParseInLocation("2006-01-02", invoice.Date, loc)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	date = date.UTC()
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}
	bill.Date = date
	bill.Where = invoice.Where
	bill.Total = 0

	for _, line := range invoice.InvoiceLines {
		description, content, unit := extractContent(line.Description)

		bi := models.BillItem{
			Description: description,
			Quantity:    line.Quantity,
			UnitValue:   line.UnitValue,
			Discount:    line.Discount / 100,
			Content:     float32(content),
			Unit:        unit,
		}
		bill.Items = append(bill.Items, &bi)
	}
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
	return b.billRepository.InsertBill(bill)
}

func extractContent(input string) (string, float64, string) {
	re := regexp.MustCompile(`(?i)(\d+)\s*(UND|u|ML|GR|g)`)
	matches := re.FindStringSubmatch(input)
	content := 0.0
	unit := ""
	if len(matches) == 3 {
		content, _ = strconv.ParseFloat(matches[1], 32)
		unit = matches[2]
	}
	switch strings.ToUpper(unit) {
	case "UND":
		unit = "UN"
	case "ML":
		unit = "ML"
	case "GR", "G":
		if content <= 1000 {
			content /= 1000
		}
		unit = "KG"
	default:
		unit = ""
	}
	description := re.ReplaceAllString(input, "")
	description = strings.TrimSpace(description)
	description = strings.ToLower(description)
	if len(description) > 0 {
		description = strings.ToUpper(string(description[0])) + description[1:]
	}
	if strings.HasSuffix(description, "*") || strings.HasSuffix(description, "x") {
		description = strings.TrimSuffix(description, "*")
		description = strings.TrimSuffix(description, "x")
	}
	return description, content, unit
}

func extractCDATA(input string) string {
	re := regexp.MustCompile(`<!\[CDATA\[(.*?)\]\]>`)
	matches := re.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}
	return input
}
