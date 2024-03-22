package repository

import (
	"github.com/raulito1500/merkadapp/src/models"
)

type MarketListRepository interface {
	ListMarketList(id string) (models.MarketList, error)
	ListMarketLists() []*models.MarketList
	InsertMarketList(marketList *models.MarketList) (string, error)
	SuggestMarketList() models.MarketList
	MarkItemCheck(idMarketList string, idProduct string) error
}
