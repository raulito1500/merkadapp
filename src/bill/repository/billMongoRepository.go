package repository

import (
	"context"
	"time"

	"github.com/raulito1500/merkadapp/src/bill/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (m *BillMongoRepository) MarkSpentItem(idBill string, idItem string) error {
	objId, _ := primitive.ObjectIDFromHex(idBill)
	filter := bson.D{
		{"_id", objId},
		{"items", bson.M{"$elemMatch": bson.M{"_id": idItem}}},
	}

	update := bson.D{
		{"$set", bson.M{"items.$.spent_date": time.Now()}},
	}
	_, err := m.coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}
	return nil
}
