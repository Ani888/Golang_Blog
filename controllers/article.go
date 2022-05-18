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

var articleCollection *mongo.Collection = configs.GetCollection(configs.DB, "articles")
var validate = validator.New()

var likeCollection *mongo.Collection = configs.GetCollection(configs.DB, "likes")
var likevalidate = validator.New()

var unlikeCollection *mongo.Collection = configs.GetCollection(configs.DB, "likes")
var unlikevalidate = validator.New()

func PostArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var article models.Article
		defer cancel()

		if err := c.BindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&article); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newArticle := models.Article{
			Id:      primitive.NewObjectID(),
			Title:   article.Title,
			Content: article.Content,
		}

		result, err := articleCollection.InsertOne(ctx, newArticle)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
	}
}

func GetArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var article models.Article
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		err := articleCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&article)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": article}})
	}
}

func GetallArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var articles []models.Article
		defer cancel()

		results, err := articleCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Article
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			articles = append(articles, singleUser)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": articles}})
	}
}

func UpdateArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var article models.Article
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&article); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if validationErr := validate.Struct(&article); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		update := bson.M{"title": article.Title, "content": article.Content}
		result, err := articleCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedArticle models.Article
		if result.MatchedCount == 1 {
			err := articleCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedArticle)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedArticle}})
	}
}

func DeleteArticle() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		result, err := articleCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.Response{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found"}})
			return
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "user successfully deleted!"}})
	}
}

// func Like() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		articleId := c.Param("articleId")
// 		defer cancel()

// 		objId, _ := primitive.ObjectIDFromHex(articleId)

// 	}
// }

func PostLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var like models.Like
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&like); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "bindingerror", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if likeValidationErr := likevalidate.Struct(&like); likeValidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "validationerror", Data: map[string]interface{}{"data": likeValidationErr.Error()}})
			return
		}

		totalLike := models.Like{
			Likecount: like.Likecount,
			ArticleId: objId,
		}

		result, err := likeCollection.InsertOne(ctx, totalLike)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})

	}

}

func GetallLikes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var likes []models.Like
		defer cancel()

		results, err := likeCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Like
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			likes = append(likes, singleUser)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": likes}})
	}
}

// func GetallLikes() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 		articleId := c.Param("articleId")
// 		var likes models.Like
// 		defer cancel()

// 		objId, _ := primitive.ObjectIDFromHex(articleId)

// 		results, err := likeCollection.Find(ctx, bson.M{})

// 		if err != nil {
// 			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 			return
// 		}

// 		defer results.Close(ctx)
// 		for results.Next(ctx) {
// 			if err = results.Decode(&likes); err != nil {
// 				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
// 			}

// 		}

// 		totalLike := models.Like{
// 			Likecount: likes.Likecount,
// 			ArticleId: objId,
// 		}

// 		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": totalLike}})
// 	}
// }

func UpdateLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var like models.Like
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&like); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if likevalidationErr := likevalidate.Struct(&like); likevalidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": likevalidationErr.Error()}})
			return
		}

		result, err := likeCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.D{{"$inc", bson.D{{"likecount", 1}}}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedLike models.Like
		if result.MatchedCount == 1 {
			err := likeCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedLike)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedLike}})
	}
}

func PostUnLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var unlike models.Unlike
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&unlike); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if unlikeValidationErr := unlikevalidate.Struct(&unlike); unlikeValidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": unlikeValidationErr.Error()}})
			return
		}

		totalunLike := models.Unlike{
			Unlikecount: unlike.Unlikecount,
			ArticleId:   objId,
		}

		result, err := unlikeCollection.InsertOne(ctx, totalunLike)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.Response{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})

	}

}

func GetallUnLikes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var unlikes []models.Unlike
		defer cancel()

		results, err := unlikeCollection.Find(ctx, bson.M{})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.Unlike
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			}

			unlikes = append(unlikes, singleUser)
		}
		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": unlikes}})
	}
}

func UpdateUnLike() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		articleId := c.Param("articleId")
		var unlike models.Unlike
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(articleId)

		if err := c.BindJSON(&unlike); err != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "binderror", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		if unlikevalidationErr := unlikevalidate.Struct(&unlike); unlikevalidationErr != nil {
			c.JSON(http.StatusBadRequest, responses.Response{Status: http.StatusBadRequest, Message: "validateerror", Data: map[string]interface{}{"data": unlikevalidationErr.Error()}})
			return
		}
		// update := bson.M{"unlikecount": unlike.Unlikecount}
		result, err := unlikeCollection.UpdateMany(ctx, bson.M{"id": objId}, bson.D{{"$inc", bson.D{{"unlikecount", 1}}}})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		var updatedunLike models.Unlike
		if result.MatchedCount == 1 {
			err := unlikeCollection.FindOneAndUpdate(ctx, bson.M{"id": objId}, bson.D{{"$inc", bson.D{{"unlikecount", 1}}}}).Decode(&updatedunLike)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.Response{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
				return
			}
		}

		c.JSON(http.StatusOK, responses.Response{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedunLike}})
	}
}
