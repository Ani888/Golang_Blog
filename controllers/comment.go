package controllers

import (
	"context"
	"gin-mongo-api-article/configs"
	"gin-mongo-api-article/models"
	"gin-mongo-api-article/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var commentCollection *mongo.Collection = configs.GetCollection(configs.DB, "comments")
var commentValidate = validator.New()

func PostComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var comment models.Comment
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if commentValidationErr := commentValidate.Struct(&comment); commentValidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": commentValidationErr.Error()}})
			return
		}

		newComment := models.Comment{
			Id:        primitive.NewObjectID(),
			Comments:  comment.Comments,
			ArticleId: objId,
		}

		result, err := commentCollection.InsertOne(ctx, newComment)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})

	}

}

func GetAllComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var commentss []models.Comment
		defer cancel()

		results, err := commentCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singlecomment models.Comment
			if err = results.Decode(&singlecomment); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			commentss = append(commentss, singlecomment)

		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": commentss}})
	}
}

func UpdateComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		commentId := c.Param("commentId")
		var comment models.Comment
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(commentId)

		if err := c.BindJSON(&comment); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if validationErr := validate.Struct(&comment); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"comments": comment.Comments}
		result, err := commentCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedComment models.Comment
		if result.MatchedCount == 1 {
			err := commentCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedComment)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedComment}})
	}
}

func DeleteComment() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		commentId := c.Param("commentId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(commentId)

		result, err := commentCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}})
			return
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}})
	}
}
