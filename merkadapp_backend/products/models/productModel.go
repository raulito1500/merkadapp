package models

var CATEGORIES = []string{"CANNED", "DELI", "PASTA", "CLEANERS", "FRUITS", "VEGETABLES", "SAUCES", "BEVERAGE", "DAIRY", "FROZEN", "PERSONAL_CARE", "SNACKS", "MEAT", "CONDIMENTS", "BAKERY", "SEAFOOD", "UNCATEGORIZED"}

type Product struct {
	ID       string `json:"id,omitempty" bson:"_id,omitempty"`
	Category string `json:"category" bson:"category"`
	Name     string `json:"name" bson:"name"`
	Quantity uint16 `json:"quantity" bson:"quantity"`
	IsBase   bool   `json:"is_base" bson:"is_base"`
	Repeat   string `json:"repeat" bson:"repeat"`
	RepeatMS int    `json:"repeatms,omitempty" bson:"repeatms"`
}

func (p *Product) GenMS() {
	ms := 0
	switch p.Repeat {
	case "1W":
		ms = 1000 * 60 * 60 * 24 * 7
	case "2W":
		ms = 1000 * 60 * 60 * 24 * 7 * 2
	case "3W":
		ms = 1000 * 60 * 60 * 24 * 7 * 3
	case "1M":
		ms = 1000 * 60 * 60 * 24 * 30
	case "2M":
		ms = 1000 * 60 * 60 * 24 * 30 * 2
	case "3M":
		ms = 1000 * 60 * 60 * 24 * 30 * 3
	default:
		ms = 0
	}
	p.RepeatMS = ms
}
