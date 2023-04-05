package main

import (
	"example/go-store-web-service/configs"
	"example/go-store-web-service/routes"
	"net/http"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func main() {
	// router := gin.Default()
	router := gin.New()
	router.Static("/assets", "./templates/assets")
	router.LoadHTMLGlob("templates/*.html")

	router.GET("/", homePage)
	router.GET("/introduction", introduction)
	router.GET("/rest-products", restProducts)

	//run database
	configs.ConnectDB()

	//routes
	routes.ProductRoute(router)

	// router.Run("localhost:3000")
}

func homePage(c *gin.Context) {
	data := gin.H{
		"title": "MDev Fake Store API | Home Page",
	}
	c.HTML(http.StatusOK, "index.html", data)
}

func introduction(c *gin.Context) {
	data := gin.H{
		"title": "MDev Fake Store API | Introduction",
	}
	c.HTML(http.StatusOK, "introduction.html", data)
}

func restProducts(c *gin.Context) {
	data := gin.H{
		"title": "MDev Fake Store API | Rest Products",
	}
	c.HTML(http.StatusOK, "rest-products.html", data)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}
