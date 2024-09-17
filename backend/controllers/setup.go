package controllers

import (
	"github.com/BaimhonS/kab-phone/configs"
	"github.com/gofiber/fiber/v2"
)

func SetUpController(app *fiber.App, configClients configs.ConfigClients) {
	controller := app.Group("/api")
	OrderController(controller, configClients)
	UserController(controller, configClients)
	PhoneController(controller, configClients)
	CartController(controller, configClients)
}
