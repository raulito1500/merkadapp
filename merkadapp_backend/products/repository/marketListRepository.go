package repository

import (
	"github.com/raulito1500/merkadapp/products/models"
)

type MarketListRepository interface {
	InsertMarketList(marketList *models.MarketList) (string, error)
	SuggestMarketList() models.MarketList
}
