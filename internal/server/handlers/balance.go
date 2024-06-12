package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/server/utils"
)

type Withdraw struct {
	Order string `json:"order" binding:"required"`
	Sum   int64  `json:"sum"   binding:"required"`
}

type WithdrawalsReponse struct {
	Order       string                  `json:"order"        binding:"required"`
	Sum         int64                   `json:"sum"          binding:"required"`
	ProcessedAt utils.FormattedDatetime `json:"processed_at" binding:"required"`
}

type UserBalanceInfo struct {
	Current   int64 `json:"current"   binding:"required"`
	Withdrawn int64 `json:"withdrawn" binding:"required"`
}

func (s *Service) GetBalance(c *gin.Context) {
	login := c.GetString("Login")
	user, err := s.Repo.GetUser(c, login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found"})
		return
	}
	c.JSON(http.StatusOK, UserBalanceInfo{
		Current:   user.Balance,
		Withdrawn: user.Withdrawn,
	})
}

func (s *Service) WithdrawBalance(c *gin.Context) {
	login := c.GetString("Login")
	var withdraw Withdraw
	if err := c.ShouldBindJSON(&withdraw); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	tx, err := s.Repo.BeginTransaction(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback(c) //nolint:errcheck // ignore check

	user, err := tx.GetUserWithLock(c, login)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found"})
		return
	}

	if user.Balance-withdraw.Sum < 0 {
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "not enough balance"})
		return
	}

	order, err := tx.GetOrderWithLock(c, withdraw.Order)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "order is not found"})
		return
	}

	if order.Login != login {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "order is user's order"})
		return
	}

	err = tx.WithdrawUserBalance(c, login, withdraw.Sum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.WithdrawOrderBalance(c, withdraw.Order, withdraw.Sum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

func (s *Service) Withdrawls(c *gin.Context) {
	login := c.GetString("Login")

	orders, err := s.Repo.GetOrders(c, login, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, []WithdrawalsReponse{})
	}

	var withdrawalsResponse []WithdrawalsReponse
	for _, order := range orders {
		withdrawalsResponse = append(
			withdrawalsResponse,
			WithdrawalsReponse{
				Order:       order.Number,
				Sum:         order.Withdrawn,
				ProcessedAt: utils.FormattedDatetime(order.ProcessedAt),
			},
		)
	}

	c.JSON(http.StatusOK, withdrawalsResponse)
}
