package organizations

import (
	"convention.ninja/internal/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router, _ *auth.IgnoreListBuilder) {
	orgGrp := grp.Group("orgs")
	orgGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	orgGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	orgGrp.Get("/:id", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	orgGrp.Delete("/:id", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}
