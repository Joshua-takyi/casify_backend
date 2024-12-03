package controllers

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const dbName = "casify"
const collectionName = "usersAuth"

const invalidCredentials = "invalid credentials"

var client *mongo.Client

func ConnectToMongoDB() error {
	// Load environment variables
	if err := godotenv.Load(".env.local"); err != nil {
		return fmt.Errorf("error loading .env file: %v", err)
	}

	// Get connection details
	uri := os.Getenv("MONGODB_URI")
	pass := os.Getenv("MONGODB_PASS")

	// Replace password placeholder
	fullUri := strings.Replace(uri, "<db_password>", pass, 1)

	// Set client options
	clientOptions := options.Client().ApplyURI(fullUri)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	var err error
	client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %v", err)
	}

	// Verify the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to ping MongoDB: %v", err)
	}

	fmt.Println("Connected to MongoDB successfully")
	return nil
}

func DisconnectFromMongoDB() error {
	if client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		return fmt.Errorf("error disconnecting from MongoDB: %v", err)
	}

	fmt.Println("Disconnected from MongoDB")
	return nil
}

// Usage example in main or another initialization function
func init() {
	if err := ConnectToMongoDB(); err != nil {
		log.Fatalf("MongoDB connection failed: %v", err)
	}
}
