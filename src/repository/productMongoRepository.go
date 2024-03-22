package repository

import (
	"context"
	"fmt"

	"github.com/raulito1500/merkadapp/src/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Conecta a la BD
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

func (m *ProductMongoRepository) ListProducts() []*models.Product {
	cursor, err := m.coll.Find(context.TODO(), bson.D{{}})

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found\n")
		return nil
	}
	if err != nil {
		panic(err)
	}

	var results []*models.Product
	if err = cursor.All(context.TODO(), &results); err != nil {
		panic(err)
	}
	return results
}

func (m *ProductMongoRepository) InsertProduct(product *models.Product) (string, error){
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
