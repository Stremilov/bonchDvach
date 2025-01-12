package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary      Add new user
// @Description  add user to the database
// @Tags         users
// @Accept       json
// @Produce      json
// @Success      201
// @Router       /bonchdvach/api/users [post]
func CreateUserHandler(c *gin.Context) {
	userIP := c.ClientIP()
	query := "INSERT INTO users (ip) VALUES ($1)"
	_, err := db.Exec(query, userIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при добавлении пользвоателя в БД", "details": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
