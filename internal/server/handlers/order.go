package handlers

import (
	"io"
	"net/http"
	"time"

	"github.com/ShiraazMoollatjie/goluhn"
	"github.com/gin-gonic/gin"
)

type OrderReponse struct {
	Number     string    `json:"number"     binding:"required"`
	Status     string    `json:"status"     binding:"required"`
	Accrual    int64     `json:"accrual"    binding:"required"`
	UploadedAt time.Time `json:"uploadedAt" binding:"required"`
}

type OrderUri struct {
	Number string `uri:"number" binding:"required"`
}

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
			c.JSON(http.StatusOK, gin.H{
				"success": true,
			})
			return
		} else {
			c.JSON(http.StatusConflict, gin.H{"error": "order number is already added by other user"})
			return
		}
	}

	err = s.Repo.StoreOrder(c, number, login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"success": true,
	})
}

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
				UploadedAt: order.UploadedAt,
			},
		)
	}

	c.JSON(http.StatusOK, ordersResponse)
}

func (s *Service) GetOrder(c *gin.Context) {
	login := c.GetString("Login")

	var orderUri OrderUri
	if err := c.ShouldBindUri(&orderUri); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	order, err := s.Repo.GetOrder(c, orderUri.Number)
	if err != nil {
		c.JSON(http.StatusNoContent, gin.H{"success": true})
		return
	}

	if login != order.Login {
		c.JSON(http.StatusConflict, gin.H{"error": "order number is already added by other user"})
		return
	}

	c.JSON(http.StatusOK, OrderReponse{
		Number:     order.Number,
		Status:     order.Status,
		Accrual:    order.Accrual,
		UploadedAt: order.UploadedAt,
	})
}
