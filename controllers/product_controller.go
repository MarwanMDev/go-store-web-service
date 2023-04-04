package controllers

import (
	"context"
	"example/go-store-web-service/configs"
	"example/go-store-web-service/models"
	"example/go-store-web-service/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "products")
var validate = validator.New()

func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var product models.Product
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newProduct := models.Product{
			Id:          primitive.NewObjectID(),
			Name:        product.Name,
			Description: product.Description,
			Title:       product.Title,
		}

		result, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.ProductResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": product}})
	}
}

func EditProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(productId)

		//validate the request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&product); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"name": product.Name, "location": product.Description, "title": product.Title}
		result, err := productCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//get updated user details
		var updatedProduct models.Product
		if result.MatchedCount == 1 {
			err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedProduct}})
	}
}

func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)

		result, err := productCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.ProductResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
		)
	}
}

func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var products []models.Product
		defer cancel()

		results, err := productCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleProduct models.Product
			if err = results.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			products = append(products, singleProduct)
		}

		c.JSON(http.StatusOK,
			responses.ProductResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": products}},
		)
	}
}
