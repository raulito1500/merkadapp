package entities

import "encoding/xml"

type AttachedDocument struct {
	XMLName     xml.Name `xml:"AttachedDocument"`
	Description string   `xml:"Attachment>ExternalReference>Description"`
}

type Invoice struct {
	XMLName      xml.Name      `xml:"Invoice"`
	Date         string        `xml:"IssueDate"`
	Time         string        `xml:"IssueTime"`
	Where        string        `xml:"AccountingSupplierParty>Party>PartyName>Name"`
	InvoiceLines []InvoiceLine `xml:"InvoiceLine"`
}

type InvoiceLine struct {
	Description string  `xml:"Item>Description"`
	Quantity    float32 `xml:"InvoicedQuantity"`
	UnitValue   float32 `xml:"Price>PriceAmount"`
	Discount    float32 `xml:"AllowanceCharge>MultiplierFactorNumeric"`
}
