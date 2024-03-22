package repository

import (
	"context"

	"github.com/raulito1500/merkadapp/src/bill/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BillMongoRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

const BILL_COLLECTION = "bills"

func NewBillMongoRepository(db *mongo.Database) BillRepository {
	return &BillMongoRepository{
		db:   db,
		coll: db.Collection(BILL_COLLECTION),
	}
}

func (m *BillMongoRepository) InsertBill(bill *models.Bill) (string, error) {
	result, err := m.coll.InsertOne(context.TODO(), bill)
	if err != nil {
		return "", err
	}
	mongoId := result.InsertedID
	return mongoId.(primitive.ObjectID).Hex(), nil
}
