package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/pisarevaa/gophermart/internal/server/storage"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

type SuccessLogin struct {
	Success bool   `json:"success" binding:"required"`
	Token   string `json:"token"   binding:"required"`
}

// RegisterUser godoc
//
//	@Summary	Regiser user
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		storage.RegisterUser	true	"Body"
//	@Success	200		{object}	storage.Success			"Response"
//	@Failure	409		{object}	storage.Error			"Login is already used"
//	@Failure	500		{object}	storage.Error			"Error"
//	@Router		/api/user/register [post]
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

	token, err := utils.GenerateJWTString(s.Config.TokenExpSec, s.Config.SecretKey, user.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", token)
	c.SetCookie("token", token, int(s.Config.TokenExpSec), "/", "localhost", false, true)

	c.JSON(http.StatusOK, storage.Success{
		Success: true,
	})
}

// LoginUser godoc
//
//	@Summary	Login user
//	@Schemes
//	@Tags		Auth
//	@Accept		json
//	@Produce	json
//	@Param		request	body		storage.RegisterUser	true	"Body"
//	@Success	200		{object}	SuccessLogin			"Response"
//	@Failure	401		{object}	storage.Error			"Login is not found or password is wrong"
//	@Failure	400		{object}	storage.Error			"Incorrect request data"
//	@Failure	500		{object}	storage.Error			"Error"
//	@Router		/api/user/login [post]
func (s *Service) LoginUser(c *gin.Context) {
	var user storage.RegisterUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userInDB, err := s.Repo.GetUser(c, user.Login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "login is not found"})
		return
	}

	isCorrect, err := utils.CheckPasswordHash(user.Password, userInDB.Password, s.Config.SecretKey)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if !isCorrect {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password is wrong"})
		return
	}

	token, err := utils.GenerateJWTString(s.Config.TokenExpSec, s.Config.SecretKey, userInDB.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Header("Authorization", token)
	c.SetCookie("token", token, int(s.Config.TokenExpSec), "/", "localhost", false, true)

	c.JSON(http.StatusOK, SuccessLogin{
		Success: true,
		Token:   token,
	})
}
