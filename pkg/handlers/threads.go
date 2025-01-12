package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Thread struct {
	ID       int    `json:"id"`
	Board_id int    `json:"board_id"`
	Title    string `json:"title"`
}

type CreateThreadRequest struct {
	Board_id int    `json:"board_id"`
	Title    string `json:"title"`
}

// @Summary      Добавить новый тред
// @Description  Создает новый тред на доске и делает запись в БД
// @Tags         threads
// @Accept       json
// @Produce      json
// @Success      201
// @Router       /bonchdvach/api/threads [post]
func CreateThread(c *gin.Context) {
	var thread CreateThreadRequest

	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	query := "INSERT INTO threads VALUES (title, board_id)"
	if _, err := db.Exec(query, thread.Title, thread.Board_id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке треда в БД", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// @Summary      Получить все треды доски
// @Description  add user to the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201
// @Router       /bonchdvach/api/threads [get]
func GetAllThreads(c *gin.Context) {
	query := "SELECT * FROM threads"

	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении досок", "details": err.Error()})
		return
	}

	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		if err := rows.Scan(&t.ID, &t.Board_id, &t.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении записи о треде", "details": err.Error()})
			return
		}
		threads = append(threads, t)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": threads})
}
