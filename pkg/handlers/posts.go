package handlers

import (
	"bonchDvach/pkg/models"
	ws "bonchDvach/pkg/websockets"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PostRepository interface {
	GetAllPosts(ctx context.Context, threadID string) ([]models.Post, error)
	CreatePost(ctx context.Context, threadID string, content string) error
}

type CreatePostRequest struct {
	ThreadID string `json:"threadID" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type SuccessGettingPostsResponse struct {
	Status string        `json:"status" example:"success"`
	Posts  []models.Post `json:"posts"`
}

type PostHandler struct {
	repo  PostRepository
	wsHub *ws.Hub
}

func NewPostHandler(repo PostRepository, wsHub *ws.Hub) PostHandler {
	return PostHandler{
		repo:  repo,
		wsHub: wsHub,
	}
}

// @Summary      Добавить новый пост
// @Description  Создает новый пост, который принадлежит определенной доске и делает запись в БД. При создании нового поста отдает в вебсокет данные: "event": "post_created", "data": {"thread_id": post.ThreadID, "content": post.Content}
// @Tags         posts
// @Accept       json
// @Produce      json
// @Success      201
// @Failure 	 400 {object} BadRequestResponse "Ошибка при получении данных"
// @Failure 	 500 {object} InternalServerErrorResponse "Ошибка при вставке поста в БД"
// @Router       /bonchdvach/api/posts [post]
func (h PostHandler) CreatePost(c *gin.Context) {
	ctx := context.Background()
	var post CreatePostRequest

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	if err := h.repo.CreatePost(ctx, post.ThreadID, post.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке поста в БД", "details": err.Error()})
		return
	}

	h.wsHub.Broadcast <- gin.H{
		"event": "post_created",
		"data": gin.H{
			"thread_id": post.ThreadID,
			"content":   post.Content,
		},
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary      Получить все посты треда
// @Description  Получает все посты определенного треда
// @Tags         posts
// @Accept       json
// @Produce      json
// @Param        threadID path int true "Thread ID"
// @Success      200 {object} SuccessGettingPostsResponse "Успешное получение постов треда"
// @Failure 	 500 {object} InternalServerErrorResponse "Внутренняя ошибка"
// @Router       /bonchdvach/api/posts/{threadID} [get]
func (h PostHandler) GetAllPosts(c *gin.Context) {
	ctx := context.Background()
	threadID := c.Param("threadID")

	posts, err := h.repo.GetAllPosts(ctx, threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении постов", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": posts})
}
