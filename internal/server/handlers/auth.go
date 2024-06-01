package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

func (s *Server) RegisterUser(c *gin.Context) {
	var user storage.User
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

func (s *Server) LoginUser(c *gin.Context) {
	var user storage.User
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}
