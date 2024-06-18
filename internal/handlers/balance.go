package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/utils"
)

type Withdraw struct {
	Order string  `json:"order" binding:"required"`
	Sum   float32 `json:"sum"   binding:"required"`
}

type WithdrawalsReponse struct {
	Order       string                  `json:"order"        binding:"required"`
	Sum         float32                 `json:"sum"          binding:"required"`
	ProcessedAt utils.FormattedDatetime `json:"processed_at" binding:"required" swaggertype:"string" example:"2024-06-12T08:00:04+03:00"`
}

type UserBalanceInfo struct {
	Current   float32 `json:"current"   binding:"required"`
	Withdrawn float32 `json:"withdrawn" binding:"required"`
}

// GetBalance godoc
//
//	@Summary	Get user's balance
//	@Schemes
//	@Tags		Balance
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer"
//	@Security	ApiKeyAuth
//	@Success	200	{object}	UserBalanceInfo	"Response"
//	@Failure	401	{object}	storage.Error	"Unauthorized"
//	@Failure	500	{object}	storage.Error	"Error"
//	@Router		/api/user/balance [get]
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

// WithdrawUserBalance godoc
//
//	@Summary	Withdraw user's balance
//	@Schemes
//	@Tags		Balance
//	@Accept		json
//	@Produce	json
//	@Param		request			body	Withdraw	true	"Body"
//	@Param		Authorization	header	string		true	"Bearer"
//	@Security	ApiKeyAuth
//	@Success	200	{object}	storage.Success	"Response"
//	@Failure	401	{object}	storage.Error	"Unauthorized"
//	@Failure	402	{object}	storage.Error	"not enough balance"
//	@Failure	422	{object}	storage.Error	"Unprocessable Entity"
//	@Failure	500	{object}	storage.Error	"Error"
//	@Router		/api/user/balance/withdraw [post]
func (s *Service) WithdrawBalance(c *gin.Context) {
	login := c.GetString("Login")
	var withdraw Withdraw
	if err := c.ShouldBindJSON(&withdraw); err != nil {
		s.Logger.Info(err.Error())
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	tx, err := s.Repo.BeginTransaction(c)
	if err != nil {
		s.Logger.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer tx.Rollback(c) //nolint:errcheck // ignore check

	user, err := tx.GetUserWithLock(c, login)
	if err != nil {
		s.Logger.Info("user is not found")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user is not found"})
		return
	}

	if user.Balance-withdraw.Sum < 0 {
		s.Logger.Info("not enough balance")
		c.JSON(http.StatusPaymentRequired, gin.H{"error": "not enough balance"})
		return
	}

	order, err := tx.GetOrderWithLock(c, withdraw.Order)
	if err != nil {
		s.Logger.Info("order is not found")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "order is not found"})
		return
	}

	if order.Login != login {
		s.Logger.Info("order is user's order")
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "order is user's order"})
		return
	}

	err = tx.WithdrawUserBalance(c, login, withdraw.Sum)
	if err != nil {
		s.Logger.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.WithdrawOrderBalance(c, withdraw.Order, withdraw.Sum)
	if err != nil {
		s.Logger.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = tx.Commit(c)
	if err != nil {
		s.Logger.Info(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
	})
}

// GetWithdrawls godoc
//
//	@Summary	Get user's withdrawls
//	@Schemes
//	@Tags		Balance
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer"
//	@Security	ApiKeyAuth
//	@Success	200	{object}	[]WithdrawalsReponse	"Response"
//	@Success	204	{object}	[]WithdrawalsReponse	"No orders"
//	@Failure	401	{object}	storage.Error			"Unauthorized"
//	@Failure	500	{object}	storage.Error			"Error"
//	@Router		/api/user/withdrawals [get]
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
				ProcessedAt: utils.FormattedDatetime(*order.ProcessedAt),
			},
		)
	}

	c.JSON(http.StatusOK, withdrawalsResponse)
}
