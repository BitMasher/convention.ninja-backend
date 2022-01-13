package auth

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/gofiber/fiber/v2"
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
