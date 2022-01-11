package users

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/users/business"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router, ignoreList *auth.IgnoreListBuilder) {

	userGrp := grp.Group("users")
	userGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
	userGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateUserAction(c)
	})
	ignoreList.AddItem("POST", "/api/users/")
	ignoreList.AddItem("POST", "/api/users") // is this necessary?
	userGrp.Get("/:id", func(c *fiber.Ctx) error {
		return business.GetUserAction(c)
	})
	userGrp.Patch("/:id", func(c *fiber.Ctx) error {
		return business.UpdateUserAction(c)
	})
	userGrp.Delete("/:id", func(c *fiber.Ctx) error {
		return business.DeleteUserAction(c)
	})
}
