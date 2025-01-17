package models

type BillBags struct {
	Quantity float32 `json:"quantity" bson:"quantity"`
	Value    float32 `json:"value" bson:"value"`
	Total    float32 `json:"total" bson:"total"`
}