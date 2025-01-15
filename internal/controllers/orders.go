package controllers

import (
	"delivery/internal/models"
	service "delivery/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OrderController struct {
	orderService service.OrderService
}

func NewOrderController(orderService service.OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (c *OrderController) PlaceOrder(ctx *gin.Context) {
	var request struct {
		UserID    int           `json:"user_id"`
		CartItems []models.Cart `json:"cart_items"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	for i := range request.CartItems {
		request.CartItems[i].UserID = uint(request.UserID)
	}

	orderID, err := c.orderService.PlaceOrder(request.UserID, request.CartItems)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"order_id": orderID})
}
