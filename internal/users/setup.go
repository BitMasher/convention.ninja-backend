package users

import (
	"convention.ninja/internal/users/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {

	userGrp := grp.Group("users")
	userGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetUsers(c)
	})
	userGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateUser(c)
	})
	userGrp.Get("/:id", func(c *fiber.Ctx) error {
		return business.GetUser(c)
	})
	userGrp.Patch("/:id", func(c *fiber.Ctx) error {
		return business.UpdateUser(c)
	})
	userGrp.Delete("/:id", func(c *fiber.Ctx) error {
		return business.DeleteUser(c)
	})
}
