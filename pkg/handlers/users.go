package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserRepository interface {
	GetUser()
	CreateUser(ctx context.Context, userIP string) error
}

type SuccessUserResponse struct {
	Status string `json:"status" example:"success"`
}

type UserHandler struct {
	repo UserRepository
}

func NewUserHandler(repo UserRepository) UserHandler {
	return UserHandler{
		repo: repo,
	}
}

// @Summary      Add new user
// @Description  add user to the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201 {object} SuccessUserResponse "Успешное создание пользователя"
// @Router       /bonchdvach/api/users [post]
func (h UserHandler) CreateUser(c *gin.Context) {
	ctx := context.Background()
	userIP := c.ClientIP()

	err := h.repo.CreateUser(ctx, userIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении пользвоателя в БД", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h UserHandler) GetUser(c *gin.Context) {
	panic("unimplemented")
}
