package config

import (
	"context"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/v2/event"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

var (
	clientInstance *mongo.Client
	once           sync.Once // Ensures connection only happens once
)

// Direct the driver to the specific host ports you mapped
const MONGODB_URI = "mongodb://mongo1:27017,mongo2:27017,mongo3:27017/?replicaSet=myReplicaSet&serverSelectionTimeoutMS=2000" // replica name should be same as in docker-compose

func GetMongoClient() (*mongo.Client, error) {
	var err error
	once.Do(func() {
		monitor := &event.CommandMonitor{
			Started: func(ctx context.Context, e *event.CommandStartedEvent) {
				// ConnectionID contains the address (e.g., mongo2:27017)
				log.Printf("MONGODB COMMAND: %s | NODE: %s | DB: %s",
					e.CommandName,
					e.ConnectionID,
					e.DatabaseName,
				)
			},
		}
		clientOptions := options.Client().
			ApplyURI(MONGODB_URI).
			SetReadPreference(readpref.SecondaryPreferred()).
			SetMonitor(monitor)

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		clientInstance, err = mongo.Connect(clientOptions)
		if err == nil {
			err = clientInstance.Ping(ctx, readpref.Primary())
		}
	})
	return clientInstance, err
}
