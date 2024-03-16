package services

import (
	"github.com/raulito1500/merkadapp/products/models"
	"github.com/raulito1500/merkadapp/products/repository"
)

type MarketListService struct {
	marketListRepository repository.MarketListRepository
}

func NewMarketListService(r repository.MarketListRepository) MarketListService {
	return MarketListService{
		marketListRepository: r,
	}
}
func (b *MarketListService) InsertMarketList(marketList *models.MarketList) (string, error) {

	return b.marketListRepository.InsertMarketList(marketList)
}

func (b *MarketListService) SuggestMarketList() models.MarketList{
	return b.marketListRepository.SuggestMarketList()
}
