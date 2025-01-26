package handlers

import (
	"bonchDvach/pkg/models"
	ws "bonchDvach/pkg/websockets"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BoardRepository interface {
	CreateBoard(ctx context.Context, name string, description string) error
	GetBoards(ctx context.Context) ([]models.Board, error)
}

type CreateBoardRequest struct {
	Name        string `json:"name" example:"Мотоциклы" binding:"required"`
	Description string `json:"description" example:"Обсуждение Питерских мотосходок" binding:"required"`
}

type SuccessGettingBoardsResponse struct {
	Status string         `json:"status" example:"success"`
	Boards []models.Board `json:"boards"`
}

type BoardHandler struct {
	repo  BoardRepository
	wsHub *ws.Hub
}

func NewBoardHandler(repository BoardRepository, wsHub *ws.Hub) BoardHandler {
	return BoardHandler{
		repo:  repository,
		wsHub: wsHub,
	}
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
func (h BoardHandler) CreateBoard(c *gin.Context) {
	var BoardRequest CreateBoardRequest
	ctx := c.Request.Context()

	if err := c.ShouldBindJSON(&BoardRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при получении данных", "details": err.Error()})
		return
	}

	err := h.repo.CreateBoard(ctx, BoardRequest.Name, BoardRequest.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при создании записи о доске в БД", "details": err.Error()})
		return
	}

	h.wsHub.Broadcast <- gin.H{
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
func (h BoardHandler) GetBoards(c *gin.Context) {
	ctx := context.Background()

	boards, err := h.repo.GetBoards(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении досок", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": boards})
}
