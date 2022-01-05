package main

import (
	"convention.ninja/internal/events"
	"convention.ninja/internal/inventory"
	"convention.ninja/internal/organizations"
	fiber "github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	apiGrp := app.Group("api/")
	organizations.SetupRoutes(apiGrp)
	inventory.SetupRoutes(apiGrp)
	events.SetupRoutes(apiGrp)
}
