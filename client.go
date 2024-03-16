package friendlymongo

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoClient is a struct to manage the database connection
type MongoClient struct {
	client      *mongo.Client
	initialized bool
}

var (
	instance *MongoClient
	once     sync.Once
)

// GetInstance returns the singleton instance of DatabaseConnection
func GetInstance() *MongoClient {
	return instance
}

// SetInstance initializes a new database connection
func SetInstance(uri string) *MongoClient {
	once.Do(func() {

		clientOptions := options.
			Client().
			ApplyURI(uri)

		c, err := mongo.Connect(context.Background(), clientOptions)
		if err != nil {
			panic(err)
		}

		instance = &MongoClient{client: c, initialized: true}
	})

	return instance
}

// Connect opens a connection to the database
func (c *MongoClient) Connect() error {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	if err := c.client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("could not connect to the database: %v", err)
	}

	return nil
}

func (c *MongoClient) Database(name string) *mongo.Database {
	if instance == nil || instance.Client() == nil {
		return nil
	}

	return c.client.Database(name)
}

func (c *MongoClient) Client() *mongo.Client {
	return c.client
}

func (c *MongoClient) Disconnect() error {
	if instance == nil || instance.Client() == nil {
		return nil
	}

	return instance.Client().Disconnect(context.Background())
}
