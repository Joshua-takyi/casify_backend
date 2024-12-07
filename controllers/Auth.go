package controllers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joshua/casify/helpers"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// RegisterClient handles the user registration process
func RegisterClient(ctx *gin.Context) {
	var inputVal model.User

	if err := ctx.ShouldBindJSON(&inputVal); err != nil {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if userExists(inputVal.Email) {

		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
			"error":   "A user with this email already exists.",
		})
		return
	}

	hashedPassword, err := helpers.HashPassword(inputVal.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
		return
	}
	inputVal.Password = hashedPassword // Update the password with the hashed one

	inputVal.Id = primitive.NewObjectID()
	inputVal.TimeStamp.CreatedAt = time.Now()
	inputVal.TimeStamp.UpdatedAt = time.Now()
	inputVal.Role = "user" // Default role is "user"

	if err := insertNewUser(inputVal); err != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to insert data",
			"error":   err.Error(),
		})
		return
	}

	if err := createUserCollection(inputVal.Id.Hex()); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user collection",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully",
		"collection": inputVal.Id.Hex(),
	})
}

func userExists(email string) bool {

	collection := Client.Database(dbName).Collection(collectionName)
	filter := bson.M{"email": email}
	var existingUser model.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	return err == nil
}

func insertNewUser(inputVal model.User) error {
	collection := Client.Database(dbName).Collection(collectionName)
	_, err := collection.InsertOne(context.Background(), inputVal)
	return err
}

func createUserCollection(userID string) error {
	_, err := helpers.CollectionExistsOrCreate(Client, userID)
	return err
}

func LoginClient(ctx *gin.Context) {
	// Parse user input
	var inputVal model.LoginRequest
	if err := ctx.ShouldBindJSON(&inputVal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request", "error": err.Error()})
		return
	}

	// Validate input
	if inputVal.Email == "" || inputVal.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Email and password are required"})
		return
	}

	// Check if user exists
	collection := Client.Database(dbName).Collection(collectionName)
	var user model.User
	err := collection.FindOne(ctx, bson.M{"email": inputVal.Email}).Decode(&user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Verify password
	if !helpers.ComparePassword(user.Password, inputVal.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to generate token"})
		return
	}

	// Set cookie
	ctx.SetCookie(
		"Authorization", // Cookie name
		tokenString,     // Cookie value (JWT token)
		3600*24*30,      // Expiry (30 days in seconds)
		"/",             // Path
		"localhost",     // Domain (for local development)
		false,           // Secure (false for local development)
		false,           // HttpOnly (false for debugging)
	)

	// Set Authorization header as well
	ctx.Header("Authorization", "Bearer "+tokenString)

	// // Debug logging
	// fmt.Printf("Setting cookie: Authorization=%s\n", tokenString)

	// Send success response with token in body
	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"token":  tokenString,
	})
}

func Validate(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}
	// jwt token
	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
