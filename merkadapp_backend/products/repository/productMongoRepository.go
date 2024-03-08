package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/raulito1500/merkadapp/database"
	"github.com/raulito1500/merkadapp/products/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// TODO: Conecta a la BD
type ProductMongo struct {
	db   *mongo.Database
	coll *mongo.Collection
}

const COLLECTION = "products"

func NewProductMongo() *ProductMongo {
	db := database.NewMongoDatabase().GetDb()
	return &ProductMongo{
		db:   db,
		coll: db.Collection(COLLECTION),
	}
}
func (m *ProductMongo) Close() {
	log.Println("Cerrando conexión: ", &m.db)
	if err := m.db.Client().Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

func (m *ProductMongo) ListProducts() []*models.Product {
	defer m.Close()

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

func (m *ProductMongo) UpdateProduct(id string) error {
	defer m.Close()

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
