package repository

import (
	"context"
	"errors"
	"time"

	"github.com/raulito1500/merkadapp/src/market_list/models"
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
const PRODUCT_COLLECTION = "products"

func NewMarketListMongoRepository(db *mongo.Database) MarketListRepository {
	return &MarketListMongoRepository{
		db:   db,
		coll: db.Collection(MARKET_LIST_COLLECTION),
	}
}

func (m *MarketListMongoRepository) ListMarketLists() []*models.MarketList {
	opts := options.Find().SetProjection(bson.D{{"items", 0}})
	opts = opts.SetSort(bson.D{{"date", -1}})
	opts = opts.SetLimit(5)

	cursor, err := m.coll.Find(context.TODO(), bson.D{}, opts)
	if err != nil {
		return []*models.MarketList{}
	}
	defer cursor.Close(context.TODO())

	results := []*models.MarketList{}

	if err := cursor.All(context.TODO(), &results); err != nil {
		return []*models.MarketList{}
	}
	return results
}

func (m *MarketListMongoRepository) ListMarketList(id string) (models.MarketList, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.MarketList{}, errors.New("Not found")
	}

	var result models.MarketList
	filter := bson.D{{"_id", objId}}
	err = m.coll.FindOne(context.TODO(), filter).Decode(&result)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.MarketList{}, errors.New("Not found")
		}
		return models.MarketList{}, err
	}
	return result, nil
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
	// filter := bson.D{{"repeat", "1W"}, {"is_base", true}}
	// opts := options.Find().SetProjection(bson.D{{"product_name", "$name"}, {"quantity", 1}, {"_id", 0}, {"product_id", "$_id"}})

	// cursor, err := collP.Find(context.TODO(), filter, opts)
	pipeline := bson.A{
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "bills"},
					{"let", bson.D{{"product_id", bson.D{{"$toString", "$_id"}}}}},
					{"pipeline",
						bson.A{
							bson.D{{"$unwind", bson.D{{"path", "$items"}}}},
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$eq",
													bson.A{
														"$items.product_id",
														"$$product_id",
													},
												},
											},
										},
									},
								},
							},
							bson.D{
								{"$group",
									bson.D{
										{"_id", "$items.product_id"},
										{"last_date", bson.D{{"$max", "$date"}}},
									},
								},
							},
						},
					},
					{"as", "products"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"bill_product",
						bson.D{
							{"$arrayElemAt",
								bson.A{
									"$products",
									0,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"since",
						bson.D{
							{"$subtract",
								bson.A{
									time.Now(),
									"$bill_product.last_date",
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$match",
				bson.D{
					{"$expr",
						bson.D{
							{"$or",
								bson.A{
									bson.D{
										{"$eq",
											bson.A{
												"$since",
												primitive.Null{},
											},
										},
									},
									bson.D{
										{"$gte",
											bson.A{
												"$since",
												"$repeatms",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{"$project",
				bson.D{
					{"product_name", "$name"},
					{"product_id", "$_id"},
					{"_id", 0},
					{"quantity", "$quantity"},
				},
			},
		},
	}
	cursor, err := collP.Aggregate(context.TODO(), pipeline)
	var result models.MarketList
	if err = cursor.All(context.TODO(), &result.Items); err != nil {
		panic(err)
	}
	result.Date = time.Now()
	return result
}

func (m *MarketListMongoRepository) MarkItemCheck(idMarketList string, idProduct string) error {
	ObjId, _ := primitive.ObjectIDFromHex(idMarketList)
	filter := bson.D{
		{"_id", ObjId},
		{"items", bson.M{"$elemMatch": bson.M{"product_id": idProduct}}},
	}

	update := bson.D{
		{"$set", bson.M{"items.$.checked": true}},
	}

	_, err := m.coll.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		return err
	}
	return nil
}
