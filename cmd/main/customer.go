package main

import (
	"log/slog"
	"net/http"
	"pizza-tracker-go/internal/models"

	"github.com/gin-gonic/gin"
)

type OrderFormData struct {
	PizzaTypes []string
	PizzaSizes []string
}

type OrderRequest struct {
	Name         string   `form:"name" binding:"required,min=2,max=100"`
	Phone        string   `form:"phone" binding:"required,min=10,max=15"`
	Address      string   `form:"address" binding:"required,min=10,max=200"`
	Sizes        []string `form:"size" binding:"required,min=1,dive,valid_pizza_size"`
	PizzaTypes   []string `form:"pizza_type" binding:"required,min=1,dive,valid_pizza_type"`
	Instructions []string `form:"instructions" binding:"max=200"`
}

func (h *Handler) ServeNewOrderForm(c *gin.Context) {
	c.HTML(http.StatusOK, "order.tmpl", OrderFormData{
		PizzaTypes: models.PizzaTypes,
		PizzaSizes: models.PizzaSizes,
	})
}

func (h *Handler) HandleNewOrderPost(c *gin.Context) {
	var form OrderRequest

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	orderItems := make([]models.OrderItem, len(form.Sizes))

	for i := range orderItems {
		orderItems[i] = models.OrderItem{
			Size:         form.Sizes[i],
			Pizza:        form.PizzaTypes[i],
			Instructions: form.Instructions[i],
		}
	}

	order := models.Order{
		CustomerName: form.Name,
		Phone:        form.Phone,
		Address:      form.Address,
		Status:       models.OrderStatuses[0],
		Items:        orderItems,
	}

	if err := h.orders.CreateOrder(&order); err != nil {
		slog.Error("Failed to create Order", "error", err)
		c.String(http.StatusInternalServerError, "Failed to create Order: %v", err)
		return
	}

	slog.Info("Order created successfully", "order_id", order.ID, "customer_name", order.CustomerName)

	c.Redirect(http.StatusSeeOther, "/customer/"+order.ID)

}

func (h *Handler) serveCustomer(c *gin.Context) {
	orderID := c.Param("id")
	if orderID == "" {
		c.String(http.StatusBadRequest, "Order ID is required")
	}

	order, err := h.orders.GetOrder(orderID)
	if err != nil {
		slog.Error("Failed to get order", "error", err)
		c.String(http.StatusNotFound, "Failed to get order: %v", err)
		return
	}

	c.HTML(http.StatusOK, "customer.tmpl", gin.H{
		"order": order,
	})
}
