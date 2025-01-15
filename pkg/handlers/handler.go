package handlers

import (
	_ "bonchDvach/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"time"
)

func InitRoutesAndDB() *gin.Engine {
	InitDB()
	InitWebSocketHub()

	router := gin.New()

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Status(204)
	})

	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		ExposeHeaders:    []string{"Content-Length"},
		MaxAge:           12 * time.Hour,
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
			threads.GET("/", GetAllThreads)
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
