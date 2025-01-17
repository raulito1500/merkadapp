package models

import "time"

var UNITS = []string{"KG", "ML", "UN"}
// TODO: Esto es una entity

type Bill struct {
	ID     string       `json:"id,omitempty" bson:"_id,omitempty"`
	Date   time.Time    `json:"date" bson:"date"`
	PaidBy string       `json:"paid_by" bson:"paid_by"`
	Where  string       `json:"where" bson:"where"`
	Total  float32      `json:"total" bson:"total"`
	Items  []*BillItem  `json:"items" bson:"items"`
	Bags   []*BillBags  `json:"bags" bson:"bags,omitempty"`
	Taxes  []*BillTaxes `json:"taxes" bson:"taxes,omitempty"`
}
