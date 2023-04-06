package handler

import (
	"example/go-store-web-service/controllers"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

var (
	app *gin.Engine
)

func productRoutes(r *gin.RouterGroup) {
	r.GET("/api/products", controllers.GetAllProducts())

}

// init gin app
func init() {
	app = gin.New()

	app.Static("/assets", "./templates/assets")
	app.LoadHTMLGlob("templates/*.html")

	// Handling routing errors
	app.NoRoute(func(c *gin.Context) {
		sb := &strings.Builder{}
		sb.WriteString("routing err: no route, try this:\n")
		for _, v := range app.Routes() {
			sb.WriteString(fmt.Sprintf("%s %s\n", v.Method, v.Path))
		}
		c.String(http.StatusBadRequest, sb.String())
	})

	r := app.Group("/")

	// register route
	productRoutes(r)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	app.ServeHTTP(w, r)
}

func homePage(c *gin.Context) {
	data := gin.H{
		"title": "MDev Fake Store API | Home Page",
	}
	c.HTML(http.StatusOK, "index.html", data)
}
