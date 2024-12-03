package helpers

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

// check if collection exists
func CollectionExistsOrCreate(client *mongo.Client, collectionName string) (bool, error) {
	// Access the database
	db := client.Database("casify")

	// Check if the collection exists by listing the collections
	collections, err := db.ListCollectionNames(context.TODO(), bson.D{})
	if err != nil {
		return false, fmt.Errorf("error listing collections: %v", err)
	}

	// If the collection already exists, return true
	for _, name := range collections {
		if name == collectionName {
			return true, nil
		}
	}

	// If the collection does not exist, create it explicitly
	err = db.CreateCollection(context.TODO(), collectionName)
	if err != nil {
		return false, fmt.Errorf("error creating collection: %v", err)
	}

	// Return false indicating that the collection was created
	return false, nil
}


func HashPassword(password string) (string, error) {
	bytes , err := bcrypt.GenerateFromPassword([]byte(password), 14)
	
	if err != nil{		
		log.Fatal("error hashing password", err)
	}
	
	return string(bytes), nil
}


func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
  return err == nil
}	