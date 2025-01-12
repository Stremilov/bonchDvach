package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Post struct {
	ID       int    `json:"id"`
	ThreadID int    `json:"threadID"`
	Content  string `json:"content"`
}

type CreatePostRequest struct {
	ThreadID int    `json:"threadID" binding:"required"`
	Content  string `json:"content" binding:"required"`
}

type SuccessGettingPostsResponse struct {
	Status string `json:"status" example:"success"`
	Posts  []Post `json:"posts"`
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
func CreatePost(c *gin.Context) {
	var post CreatePostRequest

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	query := "INSERT INTO posts (thread_id, content) VALUES ($1, $2)"
	if _, err := db.Exec(query, post.ThreadID, post.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке поста в БД", "details": err.Error()})
		return
	}

	wsHub.Broadcast <- gin.H{
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
func GetAllPosts(c *gin.Context) {
	threadID := c.Param("threadID")
	query := "SELECT * FROM posts WHERE thread_id = $1"

	rows, err := db.Query(query, threadID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении постов", "details": err.Error()})
		return
	}

	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.ThreadID, &p.Content); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении записи о посте", "details": err.Error()})
			return
		}
		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при завершении обработки постов", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": posts})
}
