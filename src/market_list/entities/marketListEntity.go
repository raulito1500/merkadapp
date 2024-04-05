package entities

import "time"

type MarketList struct {
	ID             string      `json:"id,omitempty" bson:"_id,omitempty"`
	Date           time.Time   `json:"date" bson:"date"`
	Items          []*ListItem `json:"items,omitempty" bson:"items"`
}

type ListItem struct {
	ID          string  `json:"id,omitempty" bson:"_id,omitempty"`
	ProductId   string  `json:"product_id" bson:"product_id"`
	ProductName string  `json:"product_name" bson:"product_name"`
	Quantity    float32 `json:"quantity" bson:"quantity"`
	Checked     bool    `json:"checked" bson:"checked"`
}
