package business

import (
	"context"
	"convention.ninja/internal/auth"
	"convention.ninja/internal/auth/guards"
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	userData "convention.ninja/internal/users/data"
	"errors"
	auth2 "firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetUsers(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func CreateUser(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return c.Status(fiber.StatusBadRequest).SendString("User already created")
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}

	if len(req.Name) == 0 || len(req.Email) == 0 || !strings.Contains(req.Email, "@") {
		return c.Status(fiber.StatusBadRequest).SendString("Bad request")
	}
	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}

	user := userData.User{
		SnowflakeModel: data.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Email:         req.Email,
		EmailVerified: false,
		FirebaseId:    c.Locals("idtoken").(auth.IdToken).UserId,
	}
	// TODO send email verification
	data.GetConn().Create(&user)
	return c.Status(fiber.StatusOK).JSON(&user)
}

func GetUser(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("id", ""), c) {
		return c.Status(fiber.StatusOK).JSON(c.Locals("user").(*userData.User))
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized request")
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func UpdateUser(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("id", ""), c) {
		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Bad request")
		}
		if len(req.Name) == 0 && len(req.DisplayName) == 0 && len(req.Email) == 0 {
			return c.Status(fiber.StatusOK).SendString("")
		}
		user := c.Locals("user").(*userData.User)
		if len(req.Name) > 0 {
			user.Name = req.Name
		}
		if len(req.DisplayName) > 0 {
			user.DisplayName = req.DisplayName
		}
		if len(req.Email) > 0 {
			user.Email = req.Email
			user.EmailVerified = false
		}
		data.GetConn().Save(&user)
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized request")
}

func DeleteUser(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("id", ""), c) {
		user := c.Locals("user").(*userData.User)
		client, err := auth.GetConfig().FirebaseApp.Auth(context.Background())
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Unexpected error")
		}
		params := (&auth2.UserToUpdate{}).Disabled(true)
		_, err = client.UpdateUser(context.Background(), c.Locals("idtoken").(*auth.IdToken).UserId, params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Unexpected error")
		}
		if data.GetConn().Delete(&user).RowsAffected > 0 {
			return c.Status(fiber.StatusOK).SendString("")
		}
	}
	return c.Status(fiber.StatusUnauthorized).SendString("Unauthorized request")
}
