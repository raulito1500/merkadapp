package entities

type BillTotal struct {
	Month uint8   `json:"month,omitempty" bson:"_id,omitempty"`
	Total float32 `json:"total,omitempty" bson:"total,truncate"`
}
