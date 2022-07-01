package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/panggggg/golang-project/configs"
	"github.com/panggggg/golang-project/models"
	"github.com/panggggg/golang-project/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var blogCollection *mongo.Collection = configs.GetCollection(configs.DB, "blogs")
var validate = validator.New()

func CreateBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var blog models.Blog
		defer cancel()

		// validate the requesr body
		if err := c.BindJSON(&blog); err != nil {
			c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return 
		}

		// use the validator libary to validate required fields
		if validatorErr := validate.Struct(&blog); validatorErr != nil {
			c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validatorErr.Error()}})
			return 
		}

		newBlog := models.Blog{
			BlogID:		primitive.NewObjectID(),
			Title: 		blog.Title,
			Details:	blog.Details,
			Author:		blog.Author,
			Created_at:	time.Now(),
		}

		result, err := blogCollection.InsertOne(ctx, newBlog)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BlogResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.BlogResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetABlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		blogId := c.Param("blogId")
		var blog models.Blog
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(blogId)

		err := blogCollection.FindOne(ctx, bson.M{"blogid": objId}).Decode(&blog)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BlogResponse{
				Status: http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{"data": err.Error()},
			})
			return
		}
		c.JSON(http.StatusOK, responses.BlogResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": blog}})
	}
}

func UpdateBlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		blogId := c.Param("blogId")
		var blog models.Blog
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(blogId)

		// validate the request body
		if err := c.BindJSON(&blog); err != nil {
			c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		// use the validator libary to validate required fields
		if validationError := validate.Struct(&blog); validationError != nil {
			c.JSON(http.StatusBadRequest, responses.BlogResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationError.Error()}})
			return
		}

		update := bson.M{"title": blog.Title, "details": blog.Details, "author": blog.Author}
		result, err := blogCollection.UpdateOne(ctx, bson.M{"blogid": objId}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BlogResponse{
				Status: http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{"data": err.Error()},
			})
			return
		}

		// get updated blog details
		var updatedBlog models.Blog
		if result.MatchedCount == 1 {
			err := blogCollection.FindOne(ctx, bson.M{"blogid": objId}).Decode(&updatedBlog)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.BlogResponse{
					Status: http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{"data": err.Error()},
				})
				return
			}
		}

		c.JSON(http.StatusOK, responses.BlogResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedBlog}})
	}
}

func DeleteABlog() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		blogId := c.Param("blogId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(blogId)

		result, err := blogCollection.DeleteOne(ctx, bson.M{"blogid": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BlogResponse{
				Status: http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{"data": err.Error()},
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.BlogResponse{
				Status: http.StatusNotFound,
				Message: "Blog not found",
				Data: map[string]interface{}{"data": err.Error()},
			})
			return
		}

		c.JSON(http.StatusOK, responses.BlogResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "Blog successfully deleted!"}})
	}
}

func GetAllBlogs() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var blogs []models.Blog
		defer cancel()

		results, err := blogCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.BlogResponse{
				Status: http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{"data": err.Error()},
			})
			return
		}

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleBlog models.Blog
			if err = results.Decode(&singleBlog); err != nil {
				c.JSON(http.StatusInternalServerError, responses.BlogResponse{
					Status: http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{"data": err.Error()},
				})
			}
			blogs = append(blogs, singleBlog)
		}

		c.JSON(http.StatusOK, responses.BlogResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": blogs}})
	}
}