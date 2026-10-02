package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/raulito1500/merkadapp/src/market_list/models"
	"github.com/raulito1500/merkadapp/src/market_list/repository"
)

type fakeRepository struct {
	repository.MarketListRepository
	ingredients []models.RecentIngredient
	err         error
}

func (f fakeRepository) RecentIngredients() ([]models.RecentIngredient, error) {
	return f.ingredients, f.err
}

func TestRecentIngredients(t *testing.T) {
	boom := errors.New("boom")
	tests := []struct {
		name     string
		repo     fakeRepository
		expected []models.RecentIngredient
		err      error
	}{
		{"returns repository ingredients", fakeRepository{ingredients: []models.RecentIngredient{{ProductId: "1", ProductName: "Rice", Category: "PASTA"}}}, []models.RecentIngredient{{ProductId: "1", ProductName: "Rice", Category: "PASTA"}}, nil},
		{"propagates repository error", fakeRepository{err: boom}, nil, boom},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewMarketListService(tt.repo)
			got, err := s.RecentIngredients()
			if err != tt.err {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
