package users

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/users/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {

	userGrp := grp.Group("users")
	userGrp.Get("/", auth.NewUserRequired(func(c *fiber.Ctx) error {
		return business.GetUsers(c)
	}))
	userGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateUser(c)
	})
	userGrp.Get("/me", auth.NewUserRequired(func(c *fiber.Ctx) error {
		return business.GetMe(c)
	}))
	userGrp.Get("/:userId", auth.NewUserRequired(func(c *fiber.Ctx) error {
		return business.GetUser(c)
	}))
	userGrp.Patch("/:userId", auth.NewUserRequired(func(c *fiber.Ctx) error {
		return business.UpdateUser(c)
	}))
	userGrp.Delete("/:userId", auth.NewUserRequired(func(c *fiber.Ctx) error {
		return business.DeleteUser(c)
	}))
}
