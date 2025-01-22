package entities

type Recommendation struct {
	ProductName string `json:"product_name,omitempty" bson:"_id,omitempty"`
	Count       int    `json:"count" bson:"count"`
}