package controllers

import (
	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/services"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"
)

func PhoneController(app fiber.Router, configClients configs.ConfigClients) {
	phoneController := app.Group("/phones")
	phoneService := services.NewPhoneService(configClients)
	userValidate := validates.NewUserValidate()
	phoneValidate := validates.NewPhoneValidate()

	phoneController.Get("", phoneService.GetPhones)
	phoneController.Get("/images/:id", phoneService.GetPhoneImageByID)
	phoneController.Post("", userValidate.ValidateRoleAdmin, phoneValidate.ValidateCreatePhone, phoneService.CreatePhone)
	phoneController.Put("/:id", userValidate.ValidateRoleAdmin, phoneValidate.ValidateUpdatePhone, phoneService.UpdatePhone)
	phoneController.Delete("/:id", userValidate.ValidateRoleAdmin, phoneService.DeletePhone)
}
