package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/syobonaction/fur_lange/controllers"
)

func PartnerRoute(app *fiber.App) {
	//All routes related to users comes here
	app.Get("/partner/:partnerId", controllers.GetPartner)
	app.Get("/partners", controllers.GetAllPartners)
}
