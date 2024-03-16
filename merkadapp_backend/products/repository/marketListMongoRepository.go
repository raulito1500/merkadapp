package repository

import (
	"context"
	"time"

	"github.com/raulito1500/merkadapp/products/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (m *MarketListMongoRepository) MarkCheck(idMarketList string, idProduct string) error {
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
