package helpers

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/joshua/casify/model"
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
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)

	if err != nil {
		log.Fatal("error hashing password", err)
	}

	return string(bytes), nil
}

func ComparePassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func ValidateProductInput(p model.Product) error {
	// Check required fields
	if p.Title == "" || p.Description == "" || p.Price == 0 ||
		p.Images == nil || p.Details.Details == nil || p.Color == "" || p.Category == nil {
		return errors.New("all fields are required")
	}

	// Validate price and discount
	if p.Price < 0 || p.Discount < 0 {
		return errors.New("price and discount must be non-negative")
	}

	// Validate ratings
	if p.Rating < 0 || p.Rating > 5 {
		return errors.New("ratings must be between 0 and 5")
	}

	// Validate comments
	if err := validateStringSlice(p.Comments, "comments"); err != nil {
		return err
	}

	// Validate images
	if err := validateStringSlice(p.Images, "images"); err != nil {
		return err
	}

	return nil
}

func validateStringSlice(slice []string, fieldName string) error {
	if len(slice) > 0 {
		for _, item := range slice {
			if item == "" {
				return fmt.Errorf("%s cannot contain empty strings", fieldName)
			}
		}
	}
	return nil
}
