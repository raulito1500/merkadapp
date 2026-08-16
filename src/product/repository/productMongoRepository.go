package repository

import (
	"context"
	"fmt"

	"github.com/raulito1500/merkadapp/src/product/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductMongoRepository struct {
	db   *mongo.Database
	coll *mongo.Collection
}

const PRODUCT_COLLECTION = "products"

func NewProductMongoRepository(db *mongo.Database) *ProductMongoRepository {
	return &ProductMongoRepository{
		db:   db,
		coll: db.Collection(PRODUCT_COLLECTION),
	}
}

func (m *ProductMongoRepository) ListProducts() []*models.ProductListItem {
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
								{"$match",
									bson.D{
										{"where", bson.D{{"$ne", ""}}},
										{"items.total", bson.D{{"$gt", 0}}},
									},
								},
							},
							bson.D{{"$sort", bson.D{{"date", -1}}}},
							bson.D{
								{"$group",
									bson.D{
										{"_id", "$items.product_id"},
										{"purchases",
											bson.D{
												{"$push",
													bson.D{
														{"date", "$date"},
														{"where", "$where"},
														{"value", "$items.total"},
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
										{"_id", 0},
										{"purchases", bson.D{{"$slice", bson.A{"$purchases", 2}}}},
									},
								},
							},
						},
					},
					{"as", "purchase_history"},
				},
			},
		},
		bson.D{
			{"$set",
				bson.D{
					{"purchases",
						bson.D{
							{"$ifNull",
								bson.A{
									bson.D{{"$arrayElemAt", bson.A{"$purchase_history.purchases", 0}}},
									bson.A{},
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
					{"_id", 1},
					{"category", 1},
					{"name", 1},
					{"quantity", 1},
					{"is_base", 1},
					{"repeat", 1},
					{"repeatms", 1},
					{"last_where", bson.D{{"$arrayElemAt", bson.A{"$purchases.where", 0}}}},
					{"last_value", bson.D{{"$arrayElemAt", bson.A{"$purchases.value", 0}}}},
					{"last_date", bson.D{{"$arrayElemAt", bson.A{"$purchases.date", 0}}}},
					{"previous_value", bson.D{{"$arrayElemAt", bson.A{"$purchases.value", 1}}}},
				},
			},
		},
		bson.D{{"$sort", bson.D{{"name", 1}}}},
	}

	opts := options.Aggregate().SetCollation(&options.Collation{Locale: "en"})
	cursor, err := m.coll.Aggregate(context.TODO(), pipeline, opts)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found\n")
		return nil
	}
	if err != nil {
		panic(err)
	}

	var results []*models.ProductListItem
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}

	for _, item := range results {
		if item.LastValue != nil && item.PreviousValue != nil && *item.PreviousValue != 0 {
			trend := (*item.LastValue - *item.PreviousValue) / *item.PreviousValue * 100
			item.TrendPercent = &trend
		}
	}

	return results
}

func (m *ProductMongoRepository) InsertProduct(product *models.Product) (string, error) {
	result, err := m.coll.InsertOne(context.TODO(), product)
	if err != nil {
		panic(err)
	}
	mongoId := result.InsertedID
	return mongoId.(primitive.ObjectID).Hex(), nil
}

// TODO: Hacer un verdadero update
func (m *ProductMongoRepository) UpdateProduct(id string) error {
	ObjId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{{Key: "_id", Value: ObjId}}
	update := bson.D{{Key: "$set", Value: bson.D{{
		Key:   "is_base",
		Value: true,
	}}}}

	_, err := m.coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		panic(err)
	}
	return nil
}
