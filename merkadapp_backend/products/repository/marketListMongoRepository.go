package repository

import (
	"context"
	"time"

	"github.com/raulito1500/merkadapp/products/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MarketListMongoRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

const MARKET_LIST_COLLECTION = "market_list"

func NewMarketListMongoRepository(db *mongo.Database) MarketListRepository {
	return &MarketListMongoRepository{
		db:   db,
		coll: db.Collection(MARKET_LIST_COLLECTION),
	}
}

func (m *MarketListMongoRepository) InsertMarketList(marketList *models.MarketList) (string, error) {
	result, err := m.coll.InsertOne(context.TODO(), marketList)
	if err != nil {
		return "", err
	}
	mongoId := result.InsertedID
	return mongoId.(primitive.ObjectID).Hex(), nil
}

func (m *MarketListMongoRepository) SuggestMarketList() models.MarketList {

	var collP = m.db.Collection(PRODUCT_COLLECTION)
	filter := bson.D{{"repeat", "1W"}, {"is_base", true}}
	opts := options.Find().SetProjection(bson.D{{"product_name", "$name"}, {"quantity", 1}, {"_id", 0}, {"product_id", "$_id"}})

	cursor, err := collP.Find(context.TODO(), filter, opts)

	var result models.MarketList
	if err = cursor.All(context.TODO(), &result.Items); err != nil {
		panic(err)
	}
	result.Date = time.Now()
	return result
}
