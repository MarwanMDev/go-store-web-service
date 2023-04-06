package routes

import (
	"example/go-store-web-service/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/user", controllers.CreateUser())
	// router.GET("/user/:productId", controllers.GetProduct())
	// router.PUT("/user/:productId", controllers.EditProduct())
	// router.DELETE("/user/:productId", controllers.DeleteProduct())
	router.GET("/users", controllers.GetAllUsers())
}
