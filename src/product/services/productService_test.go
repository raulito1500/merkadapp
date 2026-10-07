package services

import (
	"errors"
	"reflect"
	"testing"

	"github.com/raulito1500/merkadapp/src/product/models"
	"github.com/raulito1500/merkadapp/src/product/repository"
)

type fakeRepository struct {
	repository.ProductRepository
	ingredients []models.AvailableIngredient
	err         error
}

func (f fakeRepository) AvailableIngredients() ([]models.AvailableIngredient, error) {
	return f.ingredients, f.err
}

func TestAvailableIngredients(t *testing.T) {
	boom := errors.New("boom")
	rice := []models.AvailableIngredient{{ProductId: "1", ProductName: "Rice", Category: "PASTA"}}
	tests := []struct {
		name     string
		repo     fakeRepository
		expected []models.AvailableIngredient
		err      error
	}{
		{"returns repository ingredients", fakeRepository{ingredients: rice}, rice, nil},
		{"propagates repository error", fakeRepository{err: boom}, nil, boom},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewProductService(tt.repo)
			got, err := s.AvailableIngredients()
			if err != tt.err {
				t.Fatalf("err = %v, want %v", err, tt.err)
			}
			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("got %v, want %v", got, tt.expected)
			}
		})
	}
}
