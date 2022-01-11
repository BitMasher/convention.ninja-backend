package events

import (
	"convention.ninja/internal/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router, _ *auth.IgnoreListBuilder) {
	eventsGrp := grp.Group("/orgs/:orgId/events")

	eventsGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	eventsGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	eventsGrp.Get("/:eventId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	eventsGrp.Put("/:eventId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	eventsGrp.Patch("/:eventId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	calendarsGrp := eventsGrp.Group("/:eventId/calendar")

	calendarsGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	calendarsGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	calendarsGrp.Get("/:calendarId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	calendarsGrp.Put("/:calendarId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	calendarsGrp.Patch("/:calendarId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp := eventsGrp.Group("/venues")

	venuesGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Get("/:venueId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Patch("/:venueId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Put("/:venueId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Delete("/:venueId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Get("/:venueId/locations", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Put("/:venueId/locations", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	venuesGrp.Delete("/:venueId/locations/:locationId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

}
