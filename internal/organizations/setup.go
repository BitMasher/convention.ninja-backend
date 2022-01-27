package organizations

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/organizations/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	orgGrp := grp.Group("orgs")
	orgGrp.Use(auth.NewUserRequired())
	orgGrp.Get("/", business.GetOrganizations)

	orgGrp.Post("/", business.CreateOrganization)

	orgGrp.Get("/:orgId", business.GetOrganization)

	orgGrp.Patch("/:orgId", business.UpdateOrganization)

	orgGrp.Delete("/:orgId", business.DeleteOrganization)
}
