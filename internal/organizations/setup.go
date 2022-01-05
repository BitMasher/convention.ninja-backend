package organizations

import (
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	orgGrp := grp.Group("orgs")
	orgGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
	orgGrp.Get("/:id", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}
