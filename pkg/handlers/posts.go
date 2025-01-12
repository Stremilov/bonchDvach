package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Post struct {
	ID        int    `json:"id"`
	Thread_id int    `json:"thread_id"`
	Content   string `json:"content"`
}

type CreatePostRequest struct {
	Thread_id int    `json:"thread_id"`
	Content   string `json:"content"`
}

// @Summary      Добавить новый пост
// @Description  Создает новый пост на доске и делает запись в БД
// @Tags         posts
// @Accept       json
// @Produce      json
// @Success      201
// @Router       /bonchdvach/api/posts [post]
func CreatePost(c *gin.Context) {
	var post CreatePostRequest

	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	query := "INSERT INTO posts VALUES (thread_id, content)"
	if _, err := db.Exec(query, post.Thread_id, post.Content); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке поста в БД", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary      Получить все посты треда
// @Description  Получает все посты определенного треда
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201
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
		if err := rows.Scan(&p.ID, &p.Thread_id, &p.Content); err != nil {
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
