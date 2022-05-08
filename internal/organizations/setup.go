package organizations

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/organizations/business"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	orgGrp := grp.Group("orgs")
	orgGrp.Use(auth.NewUserRequired())
	setupOrganizations(orgGrp)

	venuesGrp := orgGrp.Group("/:orgId/venues")

	setupVenues(venuesGrp)
}

func setupOrganizations(orgGrp fiber.Router) {
	orgGrp.Get("/", business.GetOrganizations)

	orgGrp.Post("/", business.CreateOrganization)

	orgGrp.Get("/:orgId", business.GetOrganization)

	orgGrp.Patch("/:orgId", business.UpdateOrganization)

	orgGrp.Delete("/:orgId", business.DeleteOrganization)
}

func setupVenues(venuesGrp fiber.Router) {
	venuesGrp.Get("/", business.GetVenues)

	venuesGrp.Post("/", business.CreateVenue)

	venuesGrp.Get("/:venueId", business.GetVenue)

	venuesGrp.Patch("/:venueId", business.UpdateVenue)

	venuesGrp.Delete("/:venueId", business.DeleteVenue)

	venuesGrp.Get("/:venueId/rooms", business.GetVenueRooms)

	venuesGrp.Post("/:venueId/rooms", business.CreateVenueRoom)

	venuesGrp.Delete("/:venueId/rooms/:roomId", business.DeleteVenueRoom)
}
