package models

var CATEGORIES = []string{"CANNED", "DELI", "PASTA", "CLEANERS", "FRUITS", "VEGETABLES", "SAUCES", "BEVERAGE", "DAIRY", "FROZEN", "PERSONAL_CARE", "SNACKS", "MEAT", "CONDIMENTS", "BAKERY", "SEAFOOD", "UNCATEGORIZED"}

type Product struct {
	ID       string `json:"id" bson:"_id,omitempty"`
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	Quantity uint16 `json:"quantity" bson:"quantity"`
	IsBase   bool   `json:"is_base" bson:"is_base"`
	Repeat   string `json:"repeat" bson:"repeat"`
}
