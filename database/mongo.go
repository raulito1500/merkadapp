package database

import (
	"context"
	"log"

	"github.com/raulito1500/merkadapp/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDatabase struct {
	Db *mongo.Database
}

func NewMongoDatabase(config *config.Config) *MongoDatabase {

	opt := options.Client().ApplyURI(config.DatabaseUrl)
	client, err := mongo.Connect(context.TODO(), opt)
	if err != nil {
		panic(err)
	}
	return &MongoDatabase{
		Db: client.Database(config.DatabaseName),
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
