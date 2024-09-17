package controllers

import (
	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/services"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"
)

func CartController(app fiber.Router, configClients configs.ConfigClients) {
	cartController := app.Group("/carts")
	cartService := services.NewCartService(configClients)
	cartValidate := validates.NewCartValidate()

	cartController.Get("/", cartService.GetCart)
	cartController.Post("/items", cartValidate.ValidateAddItemToCart, cartService.AddItemToCart)
	cartController.Delete("/items/:id", cartService.RemoveItemFromCart)
	cartController.Patch("/items/:id", cartValidate.ValidateUpdateitemFromCart, cartService.UpdateItemFromCart)
}
