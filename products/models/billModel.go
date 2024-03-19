package models

import "time"

var UNITS = []string{"KG", "ML", "UN"}

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

type BillItem struct {
	ID               string    `json:"id,omitempty" bson:"_id,omitempty"`
	ProductId        string    `json:"product_id" bson:"product_id"`
	Description      string    `json:"description" bson:"description"`
	Brand            string    `json:"brand" bson:"brand"`
	Quantity         float32   `json:"quantity" bson:"quantity"`
	Content          float32   `json:"content" bson:"content"`
	Unit             string    `json:"unit" bson:"unit"`
	IsAdditional     bool      `json:"is_additional" bson:"is_additional"`
	UnitValue        float32   `json:"unit_value" bson:"unit_value"`
	Discount         float32   `json:"discount" bson:"discount"`
	Total            float32   `json:"total" bson:"total"`
	StartConsumeDate time.Time `json:"start_consume_date" bson:"start_consume_date"`
	SpentDate        time.Time `json:"spent_date" bson:"spent_date"`
}

type BillBags struct {
	Quantity float32 `json:"quantity" bson:"quantity"`
	Value    float32 `json:"value" bson:"value"`
	Total    float32 `json:"total" bson:"total"`
}

type BillTaxes struct {
	Concept string  `json:"concept" bson:"concept"`
	Total   float32 `json:"total" bson:"total"`
}
