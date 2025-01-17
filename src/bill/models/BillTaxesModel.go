package models

type BillTaxes struct {
	Concept string  `json:"concept" bson:"concept"`
	Total   float32 `json:"total" bson:"total"`
}