package organizations

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/organizations/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	orgGrp := grp.Group("orgs")
	orgGrp.Use(auth.NewUserRequired())
	orgGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetOrganizations(c)
	})

	orgGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateOrganization(c)
	})

	orgGrp.Get("/:orgId", func(c *fiber.Ctx) error {
		return business.GetOrganization(c)
	})

	orgGrp.Patch("/:orgId", func(c *fiber.Ctx) error {
		return business.UpdateOrganization(c)
	})

	orgGrp.Delete("/:orgId", func(c *fiber.Ctx) error {
		return business.DeleteOrganization(c)
	})
}
