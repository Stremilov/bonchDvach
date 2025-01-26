package handlers

import (
	"bonchDvach/pkg/models"
	ws "bonchDvach/pkg/websockets"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ThreadRepository interface {
	GetAllThreads(ctx context.Context, boardID string) ([]models.Thread, error)
	CreateThread(ctx context.Context, title string, boardID string) error
}

type CreateThreadRequest struct {
	BoardID string `json:"boardID" binding:"required"`
	Title   string `json:"title" binding:"required"`
}

type SuccessCreateThreadResponse struct {
	Success string `json:"success" example:"success"`
}

type SuccessGetThreadsResponse struct {
	Success string          `json:"success" example:"success"`
	Threads []models.Thread `json:"threads"`
}

type ThreadsHandler struct {
	repo  ThreadRepository
	wsHub *ws.Hub
}

func NewThreadHandler(repo ThreadRepository, wsHub *ws.Hub) ThreadsHandler {
	return ThreadsHandler{
		repo:  repo,
		wsHub: wsHub,
	}
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
func (h ThreadsHandler) CreateThread(c *gin.Context) {
	ctx := c.Request.Context()
	var thread CreateThreadRequest

	if err := c.ShouldBindJSON(&thread); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Ошибка при обработке данных", "details": err.Error()})
		return
	}

	if err := h.repo.CreateThread(ctx, thread.Title, thread.BoardID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при вставке треда в БД", "details": err.Error()})
		return
	}

	h.wsHub.Broadcast <- gin.H{
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
// @Param        boardID path int true "Board ID"
// @Success      200 	{object} SuccessGetThreadsResponse "Успешное получение всех тредов"
// @Failure      500    {object}  InternalServerErrorResponse   "Внутренняя ошибка"
// @Router       /bonchdvach/api/threads/{boardID} [get]
func (h ThreadsHandler) GetAllThreads(c *gin.Context) {
	ctx := context.Background()
	boardID := c.Param("boardID")

	threads, err := h.repo.GetAllThreads(ctx, boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении тредов", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": threads})
}
