package controllers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/akhil/ecommerce-yt/database"
	"github.com/akhil/ecommerce-yt/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Application struct {
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}

func NewApplication(prodCollection, userCollection *mongo.Collection) *Application {
	return &Application{
		prodCollection: prodCollection,
		userCollection: userCollection,
	}
}

// AddToCart adds products to the cart of the user.
// GET request
// http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx
// I add dependency to this handler to show you how much simpler your handler
// become. Your handlers should validate input call one or two functions.
// These business logic functions return an answer and probably an error.
// Your handler should then handle the error and construct a repsonse that
// matches the output and error.
func (app *Application) AddToCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Validate the input of the user
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is empty")
			// Gin has this c.Errors that is supposed to be used to catch
			// all errors that are generated in your handler. I don't understand
			// why you would use that so I'm ignoring it here.
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		// Why is this called normal and not userID?
		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("user id is empty"))
			return
		}
		// Don't ignore errors!
		productID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// Do your database queries
		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		// You have sometimes multiple defer cancel() in your functions.
		defer cancel()

		err = database.AddProductToCart(ctx, app.prodCollection, app.userCollection, productID, userQueryID)
		if err != nil {
			// This error is actually controlled by us so we don't leak any
			// sensitive information about our mongodb server or what went wrong.
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Added to the cart")
	}
}

//function to remove item from cart
//GET Request
//http://localhost:8000/addtocart?id=xxxproduct_id&normal=xxxxxxuser_idxxxxxx
func (app *Application) RemoveItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		productQueryID := c.Query("id")
		if productQueryID == "" {
			log.Println("product id is inavalid")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product id is empty"))
			return
		}

		userQueryID := c.Query("userID")
		if userQueryID == "" {
			log.Println("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}

		ProductID, err := primitive.ObjectIDFromHex(productQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err = database.RemoveCartItem(ctx, app.prodCollection, app.userCollection, ProductID, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
			return
		}
		c.IndentedJSON(200, "Successfully removed from cart")
	}
}

//function to get all items in the cart and total price
//GET request
//http://localhost:8000/listcart?id=xxxxxxuser_idxxxxxxxxxx
//help to add dependency injection
//Any nice way to group together this json response?
func GetItemFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id := c.Query("id")
		if user_id == "" {
			c.Header("Content-Type", "application/json")
			c.JSON(http.StatusNotFound, gin.H{"error": "invalid id"})
			c.Abort()
			return
		}

		usert_id, _ := primitive.ObjectIDFromHex(user_id)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var filledcart models.User
		err := UserCollection.FindOne(ctx, bson.D{primitive.E{Key: "_id", Value: usert_id}}).Decode(&filledcart)
		if err != nil {
			log.Println(err)
			c.IndentedJSON(500, "not id found")
			return
		}

		filter_match := bson.D{{Key: "$match", Value: bson.D{primitive.E{Key: "_id", Value: usert_id}}}}
		unwind := bson.D{{Key: "$unwind", Value: bson.D{primitive.E{Key: "path", Value: "$usercart"}}}}
		grouping := bson.D{{Key: "$group", Value: bson.D{primitive.E{Key: "_id", Value: "$_id"}, {Key: "total", Value: bson.D{primitive.E{Key: "$sum", Value: "$usercart.price"}}}}}}
		pointcursor, err := UserCollection.Aggregate(ctx, mongo.Pipeline{filter_match, unwind, grouping})
		if err != nil {
			log.Println(err)
		}
		var listing []bson.M
		if err = pointcursor.All(ctx, &listing); err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		for _, json := range listing {
			c.IndentedJSON(200, json["total"])
			c.IndentedJSON(200, filledcart.UserCart)
		}
		ctx.Done()
	}
}

func (app *Application) BuyFromCart() gin.HandlerFunc {
	return func(c *gin.Context) {
		userQueryID := c.Query("id")
		if userQueryID == "" {
			log.Panicln("user id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		err := database.BuyItemFromCart(ctx, app.userCollection, userQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successfully Placed the order")
	}
}

func (app *Application) InstantBuy() gin.HandlerFunc {
	return func(c *gin.Context) {
		UserQueryID := c.Query("userid")
		if UserQueryID == "" {
			log.Println("UserID is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("UserID is empty"))
		}
		ProductQueryID := c.Query("pid")
		if ProductQueryID == "" {
			log.Println("Product_ID id is empty")
			_ = c.AbortWithError(http.StatusBadRequest, errors.New("product_id is empty"))
		}
		productID, err := primitive.ObjectIDFromHex(ProductQueryID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		// Always use the defer cancel() as close to the creation point as possible.
		defer cancel()
		err = database.InstantBuyer(ctx, app.prodCollection, app.userCollection, productID, UserQueryID)
		if err != nil {
			c.IndentedJSON(http.StatusInternalServerError, err)
		}
		c.IndentedJSON(200, "Successully placed the order")
	}
}
