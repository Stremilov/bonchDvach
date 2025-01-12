package handlers

import (
	_ "bonchDvach/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutesAndDB() *gin.Engine {
	InitDB()
	router := gin.New()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/bonchdvach/api")
	{
		users := api.Group("/users")
		{
			users.POST("/", CreateUserHandler)
		}
		boards := api.Group("/boards")
		{
			boards.POST("/", CreateBoard)
			boards.GET("/", GetBoards)
		}
		threads := api.Group("/threads")
		{
			threads.POST("/", CreateThread)
			threads.GET("/", GetAllThreads)
		}
		posts := api.Group("/posts")
		{
			posts.POST("/", CreatePost)
			posts.GET("/:threadID", GetAllPosts)
		}

	}

	return router
}
