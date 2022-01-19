package auth

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
	"strings"
)

type Config struct {
	FirebaseApp *firebase.App
}

type IdToken struct {
	Email  string
	UserId string
}

var config *Config

func GetConfig() *Config {
	return config
}

func New(cfg Config) fiber.Handler {
	config = &cfg
	return func(c *fiber.Ctx) error {

		// do firebase authentication
		idToken := c.Get(fiber.HeaderAuthorization)
		if len(idToken) == 0 {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired Token")
		}
		if strings.Contains(idToken, " ") {
			idToken = strings.Split(idToken, " ")[1]
		}

		client, err := cfg.FirebaseApp.Auth(context.Background())
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
		}

		token, err := client.VerifyIDTokenAndCheckRevoked(context.Background(), idToken)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired token")
		}
		c.Locals("idtoken", &IdToken{
			Email:  token.Claims["email"].(string),
			UserId: token.UID,
		})
		return c.Next()
	}
}

func NewUserRequired(next ...fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals("user") == nil {
			return c.Status(fiber.StatusUnauthorized).SendString("User registration incomplete")
		}
		if len(next) > 0 {
			return next[0](c)
		}
		return c.Next()
	}
}
