package models

import "time"

type MarketListHeader struct {
	ID             string    `json:"id,omitempty" bson:"_id,omitempty"`
	Date           time.Time `json:"date" bson:"date"`
	CompletedItems uint16    `json:"completedItems" bson:"completedItems"`
	TotalItems     uint16    `json:"totalItems" bson:"totalItems"`
	EstimatedValue float32   `json:"estimatedValue" bson:"estimatedValue,truncate"`
}

type MarketListRecommendation struct {
	ID    string                    `json:"id,omitempty" bson:"_id,omitempty"`
	Date  time.Time                 `json:"date" bson:"date"`
	Items []*ListItemRecommendation `json:"items,omitempty" bson:"items"`
}

type ListItemRecommendation struct {
	ID          string    `json:"id,omitempty" bson:"_id,omitempty"`
	ProductId   string    `json:"product_id" bson:"product_id"`
	ProductName string    `json:"product_name" bson:"product_name"`
	Quantity    float32   `json:"quantity" bson:"quantity"`
	Category    string    `json:"category" bson:"category"`
	Checked     bool      `json:"checked" bson:"checked"`
	Value       float32   `json:"value" bson:"last_value"`
	Where       string    `json:"where" bson:"last_where"`
	Date        time.Time `json:"date" bson:"last_date"`
}
