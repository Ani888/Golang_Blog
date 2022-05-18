package routes

import (
	"gin-mongo-api-article/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoute(router *gin.Engine) {
	router.POST("/article", controllers.PostArticle())
	router.GET("/article/:articleId", controllers.GetArticle())
	router.PATCH("/article/:articleId", controllers.UpdateArticle())
	router.DELETE("/article/:articleId", controllers.DeleteArticle())
	router.GET("/article", controllers.GetallArticle())
	router.POST("/article/:articleId/comment", controllers.PostComment())
	router.GET("/article/:articleId/comment", controllers.GetAllComment())
	router.PATCH("/article/:articleId/comment/:commentId", controllers.UpdateComment())
	router.DELETE("/article/:articleId/comment/:commentId", controllers.DeleteComment())
	router.POST("/article/:articleId/like", controllers.PostLike())
	router.GET("/article/:articleId/like", controllers.GetallLikes())
	router.PATCH("/article/:articleId/like", controllers.UpdateLike())
	router.POST("/article/:articleId/unlike", controllers.PostUnLike())
	router.GET("/article/:articleId/unlike", controllers.GetallUnLikes())
	router.PATCH("/article/:articleId/unlike", controllers.UpdateUnLike())
}
