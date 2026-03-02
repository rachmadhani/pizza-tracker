package main

import "github.com/gin-gonic/gin"

func setupRoutes(router *gin.Engine, h *Handler) {
	router.GET("/", h.ServeNewOrderForm)
	router.POST("/new-order", h.HandleNewOrderPost)
	router.GET("/cutomer/:id", h.serveCustomer)

	router.Static("/static", "./templates/static")
}
