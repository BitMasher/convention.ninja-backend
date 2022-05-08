package events

import (
	"convention.ninja/internal/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	eventsGrp := grp.Group("/orgs/:orgId/events")
	eventsGrp.Use(auth.NewUserRequired())
	setupEvents(eventsGrp)

	calendarsGrp := eventsGrp.Group("/:eventId/calendar")

	setupCalendars(calendarsGrp)

}

func setupCalendars(calendarsGrp fiber.Router) {
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
}

func setupEvents(eventsGrp fiber.Router) {
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
}
