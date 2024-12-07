package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joshua/casify/controllers"
	"github.com/joshua/casify/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ValidateAuth(ctx *gin.Context) {

	// get cookie of the request
	tokenString, err := ctx.Cookie("Authorization")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

	// decode/ validate the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		fmt.Printf("JWT parsing error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid token",
		})
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// check the expiration time

		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		// find the user sub
		if controllers.Client == nil {
			if err := controllers.ConnectToMongoDB(); err != nil {
				ctx.AbortWithStatus(http.StatusUnauthorized)
			}
		}

		var user model.User
		id, err := primitive.ObjectIDFromHex(claims["sub"].(string))

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "user not found",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		filter := bson.M{"_id": id}
		collection := controllers.Client.Database("casify").Collection("usersAuth")
		if err := collection.FindOne(context.Background(), filter).Decode(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"message": "user not found",
			})
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		// attach the user to the request
		type UserResponse struct {
			Id   primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
			Name string             `json:"name,omitempty" bson:"name,omitempty"`
			Role string             `json:"role,omitempty" bson:"role,omitempty"`
		}
		// we only want to return the id, name and role
		userDetails := UserResponse{
			Id:   user.Id,
			Role: user.Role,
		}
		ctx.Set("user", userDetails)
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)

	}

	ctx.Next()
}
