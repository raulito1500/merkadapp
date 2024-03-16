package services

import (
	"errors"

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
	if hasDuplicates(marketList.Items) {
		return "", errors.New("Product id duplicated")
	}
	return b.marketListRepository.InsertMarketList(marketList)
}

func (b *MarketListService) SuggestMarketList() models.MarketList {
	return b.marketListRepository.SuggestMarketList()
}

func (b *MarketListService) MarkCheck(idMarketList string, idProduct string) error {
	return b.marketListRepository.MarkCheck(idMarketList, idProduct)
}

func hasDuplicates(items []*models.ListItem) bool {
	seen := make(map[string]int)

	for _, i := range items {
		seen[i.ProductId]++
		if seen[i.ProductId] > 1 {
			return true
		}
	}
	return false
}
