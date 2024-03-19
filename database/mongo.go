package database

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Db *mongo.Database
}

func NewMongoDatabase() *MongoDatabase {
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
	return &MongoDatabase{
		Db: client.Database(DATABASE_NAME),
	}
}

func (m *MongoDatabase) GetDb() *mongo.Database {
	log.Println("Devolviendo conexión: ", &m.Db)
	return m.Db
}

// TODO: Usar esta funcion
func (m *MongoDatabase) Close() {
	log.Println("Cerrando conexión: ", &m.Db)
	if err := m.Db.Client().Disconnect(context.TODO()); err != nil {
		panic(err)
	}
}
