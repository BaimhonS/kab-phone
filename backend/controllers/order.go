package controllers

import (
	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/services"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"
)

func OrderController(app fiber.Router, configClients configs.ConfigClients) {
	orderController := app.Group("/orders")
	orderService := services.NewOrderService(configClients)
	userValidate := validates.NewUserValidate()

	orderController.Post("/confirm", orderService.ConfirmOrder)
	orderController.Get("/track-orders/:tracking_number", orderService.GetOrderByTrackingNumber)
	orderController.Get("/track-orders", orderService.GetTrackingNumbers)
	orderController.Get("/best-worst-phones", orderService.GetBestAndWorstSellingPhones)
	orderController.Get("/total-income", userValidate.ValidateRoleAdmin, orderService.GetTotalIncome)
	orderController.Get("/check-order", userValidate.ValidateRoleAdmin, orderService.GetAllOrders)
}
