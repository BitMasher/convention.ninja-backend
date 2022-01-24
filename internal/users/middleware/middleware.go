package middleware

import (
	"convention.ninja/internal/auth"
	userData "convention.ninja/internal/users/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idToken_ := c.Locals("idtoken")
		if idToken_ == nil {
			return c.Next()
		}
		if idToken, ok := idToken_.(*auth.IdToken); ok {
			if len(idToken.UserId) == 0 {
				return c.Next()
			}
			user, err := userData.GetUserByFirebase(idToken.UserId)
			if err != nil {
				fmt.Printf("got error in user middleware: %s\n", err)
				return c.Next()
			}
			if user != nil {
				c.Locals("user", &user)
			}
		}
		return c.Next()
	}
}
