package business

import (
	"context"
	"convention.ninja/internal/auth"
	"convention.ninja/internal/auth/guards"
	userData "convention.ninja/internal/users/data"
	"errors"
	auth2 "firebase.google.com/go/auth"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func GetUsers(_ *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

type CreateUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func CreateUser(c *fiber.Ctx) error {
	if c.Locals("user") != nil {
		return c.SendStatus(fiber.StatusConflict)
	}

	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if len(req.Name) == 0 || len(req.Email) == 0 || !strings.Contains(req.Email, "@") {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if len(req.DisplayName) == 0 {
		req.DisplayName = req.Name
	}

	user := userData.User{
		Name:          req.Name,
		DisplayName:   req.DisplayName,
		Email:         strings.ToLower(req.Email),
		EmailVerified: false,
		FirebaseId:    c.Locals("idtoken").(*auth.IdToken).UserId,
	}

	exists, err := userData.GetUserByFirebase(c.Locals("idtoken").(*auth.IdToken).UserId)
	if err != nil {
		fmt.Printf("got error in CreateUser: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists != nil {
		return c.Status(fiber.StatusOK).JSON(&exists)
	}
	emailExists, err := userData.EmailExists(user.Email)
	if err != nil {
		fmt.Printf("got error in CreateUser: %s\n", err) // TODO implement logging system
	}
	if emailExists {
		return c.Status(fiber.StatusConflict).SendString("email already associated with existing user")
	}

	// TODO send email verification
	err = userData.CreateUser(&user)
	if err != nil {
		fmt.Printf("got error in CreateUser: %s\n", err) // TODO implement logging system
	}
	return c.Status(fiber.StatusOK).JSON(&user)
}

func GetMe(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(c.Locals("user").(*userData.User))
}

func GetUser(c *fiber.Ctx) error {
	// TODO need to be able to get other users, but with a filtered view
	if guards.SameUserGuard(c.Params("userId", ""), c) {
		return c.Status(fiber.StatusOK).JSON(c.Locals("user").(*userData.User))
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}

type UpdateUserRequest struct {
	Name        string `json:"name"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func UpdateUser(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("userId", ""), c) {
		var req UpdateUserRequest
		if err := c.BodyParser(&req); err != nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		if len(req.Name) == 0 && len(req.DisplayName) == 0 && len(req.Email) == 0 {
			return c.SendStatus(fiber.StatusOK)
		}
		user := c.Locals("user").(*userData.User)
		if len(req.Name) > 0 {
			user.Name = req.Name
		}
		if len(req.DisplayName) > 0 {
			user.DisplayName = req.DisplayName
		}
		if len(req.Email) > 0 && strings.ToLower(req.Email) != strings.ToLower(user.Email) {
			user.Email = req.Email
			user.EmailVerified = false
		}
		if err := userData.UpdateUser(user); err != nil {
			fmt.Printf("got error in UpdateUser: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}

func DeleteUser(c *fiber.Ctx) error {
	if guards.SameUserGuard(c.Params("userId", ""), c) {
		user := c.Locals("user").(*userData.User)
		client, err := auth.GetConfig().FirebaseApp.Auth(context.Background())
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		params := (&auth2.UserToUpdate{}).Disabled(true)
		_, err = client.UpdateUser(context.Background(), c.Locals("idtoken").(*auth.IdToken).UserId, params)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if err = userData.DeleteUser(user); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		return c.SendStatus(fiber.StatusOK)
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}
