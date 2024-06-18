package handlers

import (
	"io"
	"net/http"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
	"github.com/pisarevaa/gophermart/internal/storage"
	"github.com/pisarevaa/gophermart/internal/utils"
)

type OrderReponse struct {
	Number     string                  `json:"number"     binding:"required"`
	Status     string                  `json:"status"     binding:"required"`
	Accrual    float32                 `json:"accrual"    binding:"required"`
	UploadedAt utils.FormattedDatetime `json:"uploadedAt" binding:"required" swaggertype:"string" example:"2024-06-12T08:00:04+03:00"`
}

// AddOrder godoc
//
//	@Summary	Add an order
//	@Schemes
//	@Tags		Orders
//	@Accept		plain
//	@Produce	json
//	@Param		request			body	string	true	"Body"
//	@Param		Authorization	header	string	true	"Bearer"
//	@Security	ApiKeyAuth
//	@Success	202	{object}	storage.Success	"Response"
//	@Success	200	{object}	storage.Success	"Order is already added"
//	@Failure	401	{object}	storage.Error	"Unauthorized"
//	@Failure	400	{object}	storage.Error	"Error or incorrect data"
//	@Failure	422	{object}	storage.Error	"incorrect order number"
//	@Failure	409	{object}	storage.Error	"Order number is already added by other user"
//	@Failure	500	{object}	storage.Error	"Error"
//	@Router		/api/user/orders [post]
func (s *Service) AddOrder(c *gin.Context) {
	login := c.GetString("Login")
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	number := string(body)
	if number == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order number is empty string"})
		return
	}
	err = goluhn.Validate(number)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	order, err := s.Repo.GetOrder(c, number)
	if err == nil {
		if login == order.Login {
			c.JSON(http.StatusOK, storage.Success{
				Success: true,
			})
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "order number is already added by other user"})
		}
		return
	}

	err = s.Repo.StoreOrder(c, number, login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	s.Logger.Info("successfully store order ", number, " login ", login)

	c.JSON(http.StatusAccepted, storage.Success{
		Success: true,
	})
}

// GetOrders godoc
//
//	@Summary	Get user's orders
//	@Schemes
//	@Tags		Orders
//	@Produce	json
//	@Param		Authorization	header	string	true	"Bearer"
//	@Security	ApiKeyAuth
//	@Success	200	{object}	OrderReponse	"Response"
//	@Success	204	{object}	storage.Success	"No orders"
//	@Failure	401	{object}	storage.Error	"Unauthorized"
//	@Failure	500	{object}	storage.Error	"Error"
//	@Router		/api/user/orders [get]
func (s *Service) GetOrders(c *gin.Context) {
	login := c.GetString("Login")

	orders, err := s.Repo.GetOrders(c, login, false)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(orders) == 0 {
		c.JSON(http.StatusNoContent, []OrderReponse{})
	}

	var ordersResponse []OrderReponse
	for _, order := range orders {
		ordersResponse = append(
			ordersResponse,
			OrderReponse{
				Number:     order.Number,
				Status:     order.Status,
				Accrual:    order.Accrual,
				UploadedAt: utils.FormattedDatetime(order.UploadedAt),
			},
		)
	}

	c.JSON(http.StatusOK, ordersResponse)
}
