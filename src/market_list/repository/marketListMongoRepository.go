package repository

import (
	"context"
	"errors"
	"time"

	"github.com/raulito1500/merkadapp/src/market_list/entities"
	"github.com/raulito1500/merkadapp/src/market_list/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func (m *MarketListMongoRepository) ListMarketLists() ([]*models.MarketListHeader, error) {
	pipeline := bson.A{
		bson.D{{"$sort", bson.D{{"date", -1}}}},
		bson.D{{"$limit", 5}},
		bson.D{{"$unwind", bson.D{{"path", "$items"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "bills"},
					{"let",
						bson.D{
							{"product_id", "$items.product_id"},
							{"date", "$date"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$lt",
													bson.A{
														"$date",
														"$$date",
													},
												},
											},
										},
									},
								},
							},
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
							bson.D{{"$sort", bson.D{{"date", -1}}}},
							bson.D{
								{"$group",
									bson.D{
										{"_id", "$items.product_id"},
										{"last_date", bson.D{{"$first", "$date"}}},
										{"last_where", bson.D{{"$first", "$where"}}},
										{"last_value", bson.D{{"$first", "$items.total"}}},
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
			{"$project",
				bson.D{
					{"_id", 1},
					{"date", 1},
					{"items",
						bson.D{
							{"_id", "$items._id"},
							{"product_id", "$items.product_id"},
							{"product_name", "$items.product_name"},
							{"quantity", "$items.quantity"},
							{"checked", "$items.checked"},
							{"category", "$items.category"},
							{"last_value",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_value",
											0,
										},
									},
								},
							},
							{"last_where",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_where",
											0,
										},
									},
								},
							},
							{"last_date",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_date",
											0,
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
			{"$group",
				bson.D{
					{"_id", "$_id"},
					{"estimatedValue",
						bson.D{
							{"$sum",
								bson.D{
									{"$ifNull",
										bson.A{
											bson.D{{"$toDouble",
												bson.D{
													{"$trunc",
														bson.A{
															"$items.last_value",
															2,
														},
													},
												},
											}},
											0,
										},
									},
								},
							},
						},
					},
					{"date", bson.D{{"$first", "$date"}}},
					{"completedItems",
						bson.D{
							{"$sum",
								bson.D{
									{"$cond",
										bson.D{
											{"if",
												bson.D{
													{"$eq",
														bson.A{
															"$items.checked",
															true,
														},
													},
												},
											},
											{"then", 1},
											{"else", 0},
										},
									},
								},
							},
						},
					},
					{"totalItems", bson.D{{"$sum", 1}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"date", -1}}}},
	}
	cursor, err := m.coll.Aggregate(context.TODO(), pipeline)
	if err != nil {
		return nil, errors.New("Not found")
	}
	defer cursor.Close(context.TODO())

	var results []*models.MarketListHeader

	if err := cursor.All(context.TODO(), &results); err != nil {
		return nil, errors.New("Not found")
	}

	if err := cursor.Err(); err != nil {
		return nil, errors.New("Not found")
	}
	return results, nil
}

func (m *MarketListMongoRepository) ListMarketList(id string) (models.MarketListRecommendation, error) {
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.MarketListRecommendation{}, errors.New("Not found")
	}

	pipeline := bson.A{
		bson.D{{"$match", bson.D{{"_id", objId}}}},
		bson.D{{"$unwind", bson.D{{"path", "$items"}}}},
		bson.D{
			{"$lookup",
				bson.D{
					{"from", "bills"},
					{"let",
						bson.D{
							{"product_id", "$items.product_id"},
							{"date", "$date"},
						},
					},
					{"pipeline",
						bson.A{
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$lt",
													bson.A{
														"$date",
														"$$date",
													},
												},
											},
										},
									},
								},
							},
							bson.D{{"$unwind", bson.D{{"path", "$items"}}}},
							bson.D{
								{"$match",
									bson.D{
										{"$expr",
											bson.D{
												{"$and",
													bson.A{
														bson.D{
															{"$ne",
																bson.A{
																	"$items.product_id",
																	"",
																},
															},
														},
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
									},
								},
							},
							bson.D{{"$sort", bson.D{{"date", -1}}}},
							bson.D{
								{"$group",
									bson.D{
										{"_id", "$items.product_id"},
										{"last_date", bson.D{{"$first", "$date"}}},
										{"last_where", bson.D{{"$first", "$where"}}},
										{"last_value", bson.D{{"$first", "$items.total"}}},
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
			{"$project",
				bson.D{
					{"_id", 1},
					{"date", 1},
					{"items",
						bson.D{
							{"_id", "$items._id"},
							{"product_id", "$items.product_id"},
							{"product_name", "$items.product_name"},
							{"quantity", "$items.quantity"},
							{"checked", "$items.checked"},
							{"category", "$items.category"},
							{"last_value",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_value",
											0,
										},
									},
								},
							},
							{"last_where",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_where",
											0,
										},
									},
								},
							},
							{"last_date",
								bson.D{
									{"$arrayElemAt",
										bson.A{
											"$products.last_date",
											0,
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
			{"$group",
				bson.D{
					{"_id", "$_id"},
					{"estimatedValue",
						bson.D{
							{"$sum",
								bson.D{
									{"$ifNull",
										bson.A{
											bson.D{
												{"$toDouble",
													bson.D{
														{"$trunc",
															bson.A{
																"$items.last_value",
																2,
															},
														},
													},
												},
											},
											0,
										},
									},
								},
							},
						},
					},
					{"date", bson.D{{"$first", "$date"}}},
					{"items", bson.D{{"$push", "$items"}}},
				},
			},
		},
	}
	cursor, err := m.coll.Aggregate(context.TODO(), pipeline)
	var result []models.MarketListRecommendation
	if err = cursor.All(context.TODO(), &result); err != nil {
		panic(err)
	}
	if len(result) > 0 {
		return result[0], nil
	} else {
		return models.MarketListRecommendation{}, errors.New("Not found")
	}
}
func (m *MarketListMongoRepository) InsertMarketList(marketList *entities.MarketList) (string, error) {
	result, err := m.coll.InsertOne(context.TODO(), marketList)
	if err != nil {
		return "", err
	}
	mongoId := result.InsertedID
	return mongoId.(primitive.ObjectID).Hex(), nil
}

func (m *MarketListMongoRepository) SuggestMarketList() entities.MarketList {
	SENSIBILITY := 0.88
	var collP = m.db.Collection(PRODUCT_COLLECTION)
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
					{"since",
						bson.D{
							{"$subtract",
								bson.A{
									time.Now(),
									bson.D{
										{"$arrayElemAt",
											bson.A{
												"$products.last_date",
												0,
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
			{"$set",
				bson.D{
					{"checked",
						bson.D{
							{"$ne",
								bson.A{
									"$since",
									primitive.Null{},
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
					{"checked",
						bson.D{
							{"$cond",
								bson.D{
									{"if",
										bson.D{
											{"$ne",
												bson.A{
													"$since",
													primitive.Null{},
												},
											},
										},
									},
									{"then",
										bson.D{
											{"$gte",
												bson.A{
													"$since",
													bson.D{
														{"$multiply",
															bson.A{
																"$repeatms",
																SENSIBILITY,
															},
														},
													},
												},
											},
										},
									},
									{"else", false},
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
					{"checked", "$checked"},
					{"category", "$category"},
				},
			},
		},
		bson.D{
			{"$sort",
				bson.D{
					{"category", 1},
					{"product_name", 1},
				},
			},
		},
	}
	cursor, err := collP.Aggregate(context.TODO(), pipeline)
	var result entities.MarketList
	if err = cursor.All(context.TODO(), &result.Items); err != nil {
		panic(err)
	}
	return result
}

func (m *MarketListMongoRepository) MarkItemCheck(idMarketList string, idItem string) error {
	objId, _ := primitive.ObjectIDFromHex(idMarketList)
	filter := bson.D{
		{"_id", objId},
		{"items", bson.M{"$elemMatch": bson.M{"_id": idItem}}},
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
