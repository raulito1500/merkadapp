package entities

import "time"

type BillTotal struct {
	Date time.Time `json:"date,omitempty" bson:"_id,omitempty"`
	Total float32   `json:"total" bson:"total,truncate"`
}
