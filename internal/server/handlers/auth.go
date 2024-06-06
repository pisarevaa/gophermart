package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

// PingExample godoc
//
//	@Summary	ping example
//	@Schemes
//	@Description	do ping
//	@Tags			example
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	Helloworld
//	@Router			/example/helloworld [get]
func (s *Service) RegisterUser(c *gin.Context) {
	var user storage.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := s.Repo.GetUser(c, user.Login)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "login is already used"})
		return
	}

	passwordHash, err := utils.GetPasswordHash(user.Password, s.Config.SecretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = s.Repo.StoreUser(c, user.Login, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (s *Service) LoginUser(c *gin.Context) {
	var user storage.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInDb, err := s.Repo.GetUser(c, user.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login is not found"})
		return
	}

	isCorrect, err := utils.CheckPasswordHash(user.Password, userInDb.Password, s.Config.SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !isCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password is wrong"})
		return
	}

	token, err := utils.GenerateJWTString(s.Config.TokenExpSec, s.Config.SecretKey, userInDb.Login)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"token":   token,
	})
}
