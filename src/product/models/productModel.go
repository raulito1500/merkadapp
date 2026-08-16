package models

import "time"

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

type ProductListItem struct {
	Product       `bson:",inline"`
	LastWhere     *string    `json:"last_where,omitempty" bson:"last_where,omitempty"`
	LastValue     *float32   `json:"last_value,omitempty" bson:"last_value,omitempty"`
	LastDate      *time.Time `json:"last_date,omitempty" bson:"last_date,omitempty"`
	PreviousValue *float32   `json:"-" bson:"previous_value,omitempty"`
	TrendPercent  *float32   `json:"trend_percent,omitempty" bson:"-"`
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
