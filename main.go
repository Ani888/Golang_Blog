package main

import (
	"gin-mongo-api-article/configs"
	"gin-mongo-api-article/routes"
	"github.com/gin-gonic/gin"
)

func main() {
        router := gin.Default()
		// router.GET("/", func(c *gin.Context) {
        //         c.JSON(http.StatusOK, gin.H{
        //                 "data": "Hello from Gin-gonic & mongoDB",
        //         })
        // })
		configs.ConnectDB()
		routes.UserRoute(router)
  		router.Run("localhost:6060") 
}