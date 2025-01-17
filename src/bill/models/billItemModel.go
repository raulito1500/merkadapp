package models

import "time"

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