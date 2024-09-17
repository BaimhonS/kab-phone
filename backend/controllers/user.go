package controllers

import (
	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/services"
	"github.com/BaimhonS/kab-phone/validates"
	"github.com/gofiber/fiber/v2"
)

func UserController(app fiber.Router, configClients configs.ConfigClients) {
	userController := app.Group("/users")
	userService := services.NewUserService(configClients)
	userValidate := validates.NewUserValidate()

	userController.Get("/profile", userService.GetProfileUser)
	userController.Post("/register", userValidate.ValidateRegisterUser, userService.RegisterUser)
	userController.Post("/login", userValidate.ValidateLoginUser, userService.LoginUser)
	userController.Post("/logout", userService.LogoutUser)
	userController.Put("", userValidate.ValidateUpdateUser, userService.UpdateUser)
}
