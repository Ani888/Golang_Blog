package routes

import (
	"gin-mongo-api-article/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/article", controllers.PostArticle())
	router.GET("/article/:articleId", controllers.GetArticle())
	router.GET("/article", controllers.GetallArticle())
	router.POST("/article/:articleId/comment", controllers.PostComment())
	router.GET("/article/:articleId/comment", controllers.GetAllComment())
}
