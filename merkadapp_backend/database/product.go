package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/raulito1500/merkadapp/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TODO: Conecta a la BD
type ProductMongo struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewProductMongo() *ProductMongo {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading enviroments")
	}
	DATABASE_URL := os.Getenv("DATABASE_URL")
	DATABASE_NAME := os.Getenv("DATABASE_NAME")

	opt := options.Client().ApplyURI(DATABASE_URL)
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		panic(err)
	}
	return &ProductMongo{
		client: client,
		coll:   client.Database(DATABASE_NAME).Collection(COLLECTION),
	}
}

func (m *ProductMongo) Close() {
	if err := m.client.Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}

const COLLECTION = "products"

func (m *ProductMongo) GetAll() []*models.Product {
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
