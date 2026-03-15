package cmd

import (
	"context"
	"errors"
	"log/slog"
	"os"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	DatabaseName = "chisai"
)

var db *mongo.Database
var client *mongo.Client

func InitMongoDB() error {
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return errors.New("MONGODB_URI env var not set")
	}

	c, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	client = c

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return err
	}

	db = client.Database(DatabaseName)

	slog.Info("Connected to monogdb...")

	return nil
}

func GetMongoDB() *mongo.Database {
	return db
}

func CloseMongoDB(ctx context.Context) error {
	return client.Disconnect(ctx)
}
