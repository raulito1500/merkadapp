package repository

import (
	"github.com/raulito1500/merkadapp/src/market_list/entities"
	"github.com/raulito1500/merkadapp/src/market_list/models"
)

type MarketListRepository interface {
	ListMarketList(id string) (entities.MarketList, error)
	ListMarketLists() []*models.MarketListHeader
	InsertMarketList(marketList *entities.MarketList) (string, error)
	SuggestMarketList() entities.MarketList
	MarkItemCheck(idMarketList string, idProduct string) error
}
