package models

var CATEGORIES = []string{"CANNED", "DELI", "PASTA", "CLEANERS", "FRUITS", "VEGETABLES", "SAUCES", "BEVERAGE", "DAIRY", "FROZEN", "PERSONAL_CARE", "SNACKS", "MEAT", "CONDIMENTS", "BAKERY", "SEAFOOD", "UNCATEGORIZED"}

const (
	ONE_DAY_MS   = 1000 * 60 * 60 * 24
	ONE_WEEK_MS  = ONE_DAY_MS * 7
	ONE_MONTH_MS = ONE_DAY_MS * 30
)

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
	switch p.Repeat {
	case "1W":
		p.RepeatMS = ONE_WEEK_MS
	case "2W":
		p.RepeatMS = ONE_WEEK_MS * 2
	case "3W":
		p.RepeatMS = ONE_WEEK_MS * 3
	case "1M":
		p.RepeatMS = ONE_MONTH_MS
	case "2M":
		p.RepeatMS = ONE_MONTH_MS * 2
	case "3M":
		p.RepeatMS = ONE_MONTH_MS * 3
	default:
		p.RepeatMS = 0
	}
}
