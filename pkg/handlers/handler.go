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

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Status(204)
	})

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
