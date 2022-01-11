package business

import (
	"context"
	"convention.ninja/internal/auth"
	"convention.ninja/internal/auth/guards"
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	userData "convention.ninja/internal/users/data"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func CreateUserAction(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User already created")
	}

	if auth.GetConfig() != nil {
		fb := auth.GetConfig().FirebaseApp
		client, err := fb.Auth(context.Background())
		if err != nil {
			return c.Status(fiber.StatusServiceUnavailable).SendString("Unexpected failure")
		}
		client.email
	}

	user := userData.User{
		SnowflakeModel: data.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:       c.FormValue("name"),
		Email:      c.FormValue("email"),
		FirebaseId: c.Locals("idtoken").(auth.IdToken).UserId,
	}
	data.GetConn().Create(&user)
	return c.Status(fiber.StatusOK).JSON(&user)
}

func GetUserAction(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("id", ""), c) {
		return c.Status(fiber.StatusOK).JSON(c.Locals("user").(*userData.User))
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized request")
}

func UpdateUserAction(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

func DeleteUserAction(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}
