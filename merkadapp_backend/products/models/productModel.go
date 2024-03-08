package models

type Product struct {
	ID       string `json:"id" bson:"_id"`
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	Quantity string `json:"quantity" bson:"quantity"`
	IsBase   bool   `json:"is_base" bson:"is_base"`
	Repeat   string `json:"repeat" bson:"repeat"`
}
