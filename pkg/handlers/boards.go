package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Board struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateBoardRequest struct {
	Name        string `json:"name" example:"Мотоциклы" binding:"required"`
	Description string `json:"description" example:"Обсуждение Питерских мотосходок" binding:"required"`
}

type SuccessGettingBoardsResponse struct {
	Status string  `json:"status" example:"success"`
	Boards []Board `json:"boards"`
}

// @Summary      Создать новую доску
// @Description  Создает новую доску и делает запись в БД. При создании новой доски отдает в вебсокет данные: "event": "board_created", "data": {"name": BoardRequest.Name, "description": BoardRequest.Description}
// @Tags         boards
// @Accept       json
// @Produce      json
// @Success      201
// @Param        board  body      CreateBoardRequest  true  "Информация о доске"
// @Success      201    {object}  SuccessCreatingResponse   "Успешное создание"
// @Failure      400    {object}  BadRequestResponse   "Ошибка при получении данных"
// @Failure      500    {object}  InternalServerErrorResponse   "Ошибка при создании записи о доске в БД"
// @Router       /bonchdvach/api/boards [post]
func CreateBoard(c *gin.Context) {
	var BoardRequest CreateBoardRequest

	if err := c.ShouldBindJSON(&BoardRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении данных", "details": err.Error()})
		return
	}

	query := "INSERT INTO boards (name, description) VALUES ($1, $2)"

	_, err := db.Exec(query, BoardRequest.Name, BoardRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании записи о доске в БД", "details": err.Error()})
		return
	}

	wsHub.Broadcast <- gin.H{
		"event": "board_created",
		"data": gin.H{
			"name":        BoardRequest.Name,
			"description": BoardRequest.Description,
		},
	}

	c.JSON(http.StatusCreated, gin.H{"status": "success"})
}

// @Summary      Получить все доски
// @Description  Возвращает все доски, которые есть в базе данных
// @Tags         boards
// @Accept       json
// @Produce      json
// @Success      200 {object} SuccessGettingBoardsResponse "Успешный запрос"
// @Failure      500    {object}  InternalServerErrorResponse   "Непредвиденная ошибка"
// @Router       /bonchdvach/api/boards [get]
func GetBoards(c *gin.Context) {
	query := "SELECT * FROM boards"
	rows, err := db.Query(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении досок", "details": err.Error()})
		return
	}

	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var b Board
		if err := rows.Scan(&b.ID, &b.Name, &b.Description); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении записи о доске", "details": err.Error()})
			return
		}

		boards = append(boards, b)
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": boards})
}
