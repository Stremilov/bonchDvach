package handlers

import (
	_ "bonchDvach/docs"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRoutesAndDB() *gin.Engine {
	InitDB()
	InitWebSocketHub()

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

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
			threads.GET("/:boardID", GetAllThreads)
		}
		posts := api.Group("/posts")
		{
			posts.POST("/", CreatePost)
			posts.GET("/:threadID", GetAllPosts)
		}
		ws := api.Group("/ws")
		{
			ws.GET("/", WebSocketHandler)
		}

	}

	return router
}
