package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Thread struct {
	ID      int    `json:"id"`
	BoardID int    `json:"boardID"`
	Title   string `json:"title"`
}

type CreateThreadRequest struct {
	BoardID int    `json:"boardID" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

type SuccessCreateThreadResponse struct {
	Success string `json:"success" example:"success"`
}

type SuccessGetThreadsResponse struct {
	Success string   `json:"success" example:"success"`
	Threads []Thread `json:"threads"`
}

// @Summary      Добавить новый тред
// @Description  Создает новый тред, принадлежащий определенной доске и делает запись в БД. При создании нового треда отдает в вебсокет данные: "event": "thread_created", "data": {"title": thread.Title, "board_id": thread.BoardID}
// @Tags         threads
// @Accept       json
// @Produce      json
// @Success      201 	{object} SuccessCreateThreadResponse "Успешное создание треда"
// @Failure      400    {object}  BadRequestResponse   "Ошибка при получении данных"
// @Failure      500    {object}  InternalServerErrorResponse   "Ошибка при вставке треда в БД"
// @Router       /bonchdvach/api/threads [post]
func CreateThread(c *gin.Context) {
	var thread CreateThreadRequest

	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	query := "INSERT INTO threads (title, board_id) VALUES ($1, $2)"
	if _, err := db.Exec(query, thread.Title, thread.BoardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке треда в БД", "details": err.Error()})
		return
	}

	wsHub.Broadcast <- gin.H{
		"event": "thread_created",
		"data": gin.H{
			"title":    thread.Title,
			"board_id": thread.BoardID,
		},
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

// @Summary      Получить все треды доски
// @Description  get all threads of the board
// @Tags         threads
// @Accept       json
// @Produce      json
// @Success      200 	{object} SuccessGetThreadsResponse "Успешное получение всех тредов"
// @Failure      500    {object}  InternalServerErrorResponse   "Внутренняя ошибка"
// @Router       /bonchdvach/api/threads/{boardID} [get]
func GetAllThreads(c *gin.Context) {
	query := "SELECT * FROM threads WHERE board_id = $1"
	id := c.Param("boardID")
	rows, err := db.Query(query, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении тредов", "details": err.Error()})
		return
	}

	defer rows.Close()

	var threads []Thread
	for rows.Next() {
		var t Thread
		if err := rows.Scan(&t.ID, &t.BoardID, &t.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении записи о треде", "details": err.Error()})
			return
		}
		threads = append(threads, t)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": threads})
}
