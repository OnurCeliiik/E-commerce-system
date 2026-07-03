package database

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const AuditEventsCollection = "audit_events"

type MongoDBClient struct {
	client   *mongo.Client
	database *mongo.Database
}

func ConnectFromEnv() (*MongoDBClient, error) {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		uri = "mongodb://mongo:27017"
	}

	dbName := os.Getenv("MONGODB_DATABASE")
	if dbName == "" {
		dbName = "audit_db"
	}

	return Connect(uri, dbName)
}

func Connect(uri, dbName string) (*MongoDBClient, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	db := client.Database(dbName)
	if err := ensureIndexes(ctx, db); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("ensure indexes: %w", err)
	}

	return &MongoDBClient{
		client:   client,
		database: db,
	}, nil
}

func ensureIndexes(ctx context.Context, db *mongo.Database) error {
	collection := db.Collection(AuditEventsCollection)

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "topic", Value: 1},
				{Key: "partition", Value: 1},
				{Key: "offset", Value: 1},
			},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{
				{Key: "order_id", Value: 1},
				{Key: "recorded_at", Value: 1},
			},
		},
	}

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (m *MongoDBClient) Ping(ctx context.Context) error {
	return m.client.Ping(ctx, nil)
}

func (m *MongoDBClient) Database() *mongo.Database {
	return m.database
}

func (m *MongoDBClient) Collection(name string) *mongo.Collection {
	return m.database.Collection(name)
}

func (m *MongoDBClient) Close(ctx context.Context) error {
	return m.client.Disconnect(ctx)
}
