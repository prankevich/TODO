package api

import (
	"TODO/adapter/driven/models"
	"TODO/services/contracts"
	"github.com/gin-gonic/gin"

	"net/http"
)

// AuthController — контроллер для работы с пользователями.
type AuthController struct {
	authService contracts.AuthService
}

func NewAuthController(authService contracts.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Register(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := ac.authService.Register(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "registered"})
}

func (ac *AuthController) Stats(c *gin.Context) {
	stats, err := ac.authService.Stats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}
