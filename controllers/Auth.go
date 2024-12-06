package controllers

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joshua/casify/helpers"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterClient(ctx *gin.Context) {
	inputVal := model.User{}

	if err := ctx.ShouldBindJSON(&inputVal); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
			"error":   err.Error(),
		})
		return
	}
	// check if user exists
	collection := Client.Database(dbName).Collection(collectionName)

	filter := bson.M{"email": inputVal.Email}
	var existingUser model.User
	err := collection.FindOne(context.Background(), filter).Decode(&existingUser)
	if err == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "User already exists",
			"error":   "User already exists",
		})
		return
	}
	if err != mongo.ErrNoDocuments {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to check existing user",
			"error":   err.Error(),
		})
		return
	}
	// hash password
	password := inputVal.Password
	hashedPassword, err := helpers.HashPassword(password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to hash password",
			"error":   err.Error(),
		})
		return
	}
	inputVal.Id = primitive.NewObjectID()
	inputVal.Password = hashedPassword
	inputVal.TimeStamp.CreatedAt = time.Now()
	inputVal.TimeStamp.UpdatedAt = time.Now()
	inputVal.Role = "user"

	// collection := Client.Database(dbName).Collection(collectionName)

	_, err = collection.InsertOne(context.Background(), inputVal)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to insert data",
			"error":   err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Data inserted successfully",
	})

	// create collection for the user

	_, err = helpers.CollectionExistsOrCreate(Client, inputVal.Id.Hex())

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create collection",
			"error":   err.Error(),
		})
		return
	}
	// respond success
	ctx.JSON(http.StatusCreated, gin.H{
		"message":    "Collection created successfully",
		"collection": inputVal.Id.Hex(),
	})
}

func LoginClient(ctx *gin.Context) {

	inputVal := model.LoginRequest{}

	// Bind the request body to inputVal
	if err := ctx.ShouldBindJSON(&inputVal); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"error":   "Invalid request body",
		})
		return
	}

	// Ensure email and password are provided
	if inputVal.Email == "" || inputVal.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": invalidCredentials,
			"error":   "Email and password are required",
		})
		return
	}

	// Check if user exists in the database
	collection := Client.Database(dbName).Collection(collectionName)

	var user model.User
	err := collection.FindOne(ctx, bson.M{"email": inputVal.Email}).Decode(&user)

	// If no user is found, return an error message
	if err != nil {
		if err == mongo.ErrNoDocuments {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": invalidCredentials,
				"error":   "Invalid email or password",
			})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "Authentication error",
				"error":   "An unexpected error occurred",
			})
		}
		return
	}

	// Compare the input password with the stored hashed password
	if !helpers.ComparePassword(user.Password, inputVal.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": invalidCredentials,
			"error":   "Invalid email or password",
		})
		return
	}

	//* generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.Id.Hex(),
		"exp": time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to generate token",
			"error":   err.Error(),
		})
		return
	}

	// set cookie
	ctx.SetSameSite(http.SameSiteLaxMode)

	// change cookie domain to frontend url when in production
	ctx.SetCookie("Authorization", tokenString, 3600*24*30, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, gin.H{})
}

func Validate(ctx *gin.Context) {
	user, ok := ctx.Get("user")
	if !ok {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}
