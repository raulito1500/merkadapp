package services

import (
	"errors"
	"time"

	"github.com/raulito1500/merkadapp/src/market_list/entities"
	"github.com/raulito1500/merkadapp/src/market_list/models"
	"github.com/raulito1500/merkadapp/src/market_list/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MarketListService struct {
	marketListRepository repository.MarketListRepository
}

func NewMarketListService(r repository.MarketListRepository) MarketListService {
	return MarketListService{
		marketListRepository: r,
	}
}

func (b *MarketListService) ListMarketLists() []*models.MarketListHeader {
	return b.marketListRepository.ListMarketLists()
}

func (b *MarketListService) ListMarketList(id string) (entities.MarketList, error) {
	return b.marketListRepository.ListMarketList(id)
}

func (b *MarketListService) InsertMarketList(marketList *entities.MarketList) (string, error) {
	if hasDuplicates(marketList.Items) {
		return "", errors.New("Product id duplicated")
	}
	for _, ml := range marketList.Items {
		ml.ID = primitive.NewObjectID().Hex()
	}
	return b.marketListRepository.InsertMarketList(marketList)
}

func (b *MarketListService) SuggestMarketList() entities.MarketList {
	suggested := b.marketListRepository.SuggestMarketList()
	suggested.Date = nextMarketDay()
	return suggested
}

func (b *MarketListService) MarkItemCheck(idMarketList string, idItem string) error {
	return b.marketListRepository.MarkItemCheck(idMarketList, idItem)
}

// TODO: Ubicar esto en una ubicación adecuada
func hasDuplicates(items []*entities.ListItem) bool {
	seen := make(map[string]int)
	// TODO: No tener en cuenta si el product_id está en blanco
	for _, i := range items {
		seen[i.ProductId]++
		if seen[i.ProductId] > 1 {
			return true
		}
	}
	return false
}

// TODO: Encontrar el siguiente sábado
func nextMarketDay() time.Time {
	today := time.Now()
	diasHastaSabado := 6 - int(today.Weekday())
	return today.AddDate(0, 0, diasHastaSabado)
}
