package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BaimhonS/kab-phone/configs"
	"github.com/BaimhonS/kab-phone/controllers"
	"github.com/BaimhonS/kab-phone/middlewares"
	"github.com/BaimhonS/kab-phone/scripts"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatalf("error loading .env file : %v", err)
	}
}

func main() {
	handleArgument()

	configClients := configs.SetUpConfigs()

	app := fiber.New()

	app.Use(
		cors.New(),
		logger.New(),
		middlewares.AuthToken(configClients.Redis),
	)

	controllers.SetUpController(app, configClients)

	app.Listen(fmt.Sprintf(":%s", os.Getenv("SERVER_PORT")))
}

func handleArgument() {
	isHandled := false
	for _, arg := range os.Args {
		if strings.Contains(arg, "migrate-up") {
			scripts.MigrateUp()
			isHandled = true
		}
		if strings.Contains(arg, "migrate-down") {
			scripts.MigrateDown(arg)
			isHandled = true
		}
	}

	if isHandled {
		os.Exit(0)
	}
}
