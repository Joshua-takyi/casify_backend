package router

import (
	"github.com/gin-gonic/gin"
	"github.com/joshua/casify/controllers"
)

func Router() *gin.Engine {
	r := gin.Default()
	v1 := r.Group("/api/v1")
	v1.POST("/register", controllers.RegisterClient)
	v1.POST("/login", controllers.LoginClient)
	v1.POST("/addProduct", controllers.AddProduct)
	v1.GET("/getProducts", controllers.GetProducts)
	v1.GET("/getProduct/:id", controllers.GetById)
	v1.PUT("/updateProduct/:id", controllers.UpdateProduct)
	v1.DELETE("/deleteProduct/:id", controllers.DeleteProduct)
	v1.DELETE("/deleteProducts", controllers.DeleteManyProducts)
	return r
}
