package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type MongoRepository struct {
	mongoClient *mongo.Client
}

func NewMongoRepository() *MongoRepository {
	return &MongoRepository{}
}

func (repo *MongoRepository) InsertMany(documents []interface{}, collectionName string) error {
	collection := repo.open(collectionName)
	_, err := collection.InsertMany(context.Background(), documents)
	repo.close()
	return err
}

func (repo *MongoRepository) open(collectionName string) *mongo.Collection {
	connectionString := os.Getenv("MONGO_CONNECTION_STRING")
	databaseName := os.Getenv("MONGO_DATABASE_NAME")
	clientOptions := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		fmt.Println("Error al conectar a MongoDB:", err)
		return nil
	}
	repo.mongoClient = client
	database := client.Database(databaseName)
	return database.Collection(collectionName)
}

func (repo *MongoRepository) close() {
	defer repo.mongoClient.Disconnect(context.Background())
}
