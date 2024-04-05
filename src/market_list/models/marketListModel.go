package models

import "time"

type MarketListHeader struct {
	ID             string      `json:"id,omitempty" bson:"_id,omitempty"`
	Date           time.Time   `json:"date" bson:"date"`
	CompletedItems uint16      `json:"completedItems" bson:"completedItems"`
	TotalItems     uint16      `json:"totalItems" bson:"totalItems"`
}
