package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syobonaction/fur_lange/configs"
	"github.com/syobonaction/fur_lange/routes"
)

func main() {
	app := fiber.New()

	//run database
	configs.ConnectDB()

	routes.PartnerRoute(app)

	app.Listen(":6000")
}
