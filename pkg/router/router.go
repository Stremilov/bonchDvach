package router

import (
	_ "bonchDvach/docs"
	"bonchDvach/pkg/db/postgres"
	"bonchDvach/pkg/db/postgres/repositories"
	"bonchDvach/pkg/handlers"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/joho/godotenv/autoload"
)

func InitRoutesAndDB() *gin.Engine {
	reps, err := initDB()
	if err != nil {
		log.Fatalf("cannot init db: %v", err)
	}

	bh := handlers.NewBoardHandler(reps.BoardRepository, wsHub)
	ph := handlers.NewPostHandler(reps.PostRepository, wsHub)
	th := handlers.NewThreadHandler(reps.ThreadRepository, wsHub)
	uh := handlers.NewUserHandler(reps.UserRepository)

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
			users.POST("/", uh.CreateUser)
		}
		boards := api.Group("/boards")
		{
			boards.POST("/", bh.CreateBoard)
			boards.GET("/", bh.GetBoards)
		}
		threads := api.Group("/threads")
		{
			threads.POST("/", th.CreateThread)
			threads.GET("/:boardID", th.GetAllThreads)
		}
		posts := api.Group("/posts")
		{
			posts.POST("/", ph.CreatePost)
			posts.GET("/:threadID", ph.GetAllPosts)
		}
		ws := api.Group("/ws")
		{
			ws.GET("/", WebSocketHandler)
		}

	}

	return router
}

type repostitories struct {
	BoardRepository  handlers.BoardRepository
	PostRepository   handlers.PostRepository
	ThreadRepository handlers.ThreadRepository
	UserRepository   handlers.UserRepository
}

func initDB() (reps repostitories, err error) {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSLMODE")

	if host == "" || user == "" || password == "" || dbName == "" || sslMode == "" {
		//init default values for db (will be transfered into config.y)
		host, user, password, dbName, sslMode = "localhost", "postgres", "postgres", "bonchdvach", "disable"

	}

	pool, err := postgres.New(fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=%s", host, user, password, dbName, sslMode))
	if err != nil {
		return reps, err
	}
	reps.BoardRepository = repositories.NewBoardRepository(pool)
	reps.PostRepository = repositories.NewPostRepository(pool)
	reps.ThreadRepository = repositories.NewThreadRepository(pool)
	reps.UserRepository = repositories.NewUserRepository(pool)

	return reps, nil
}
