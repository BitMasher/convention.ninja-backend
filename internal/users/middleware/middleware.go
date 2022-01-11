package middleware

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/data"
	userData "convention.ninja/internal/users/data"
	"github.com/gofiber/fiber/v2"
)

func New() fiber.Handler {
	return func(c *fiber.Ctx) error {
		idToken_ := c.Locals("idtoken")
		if idToken_ == nil {
			return c.Next()
		}
		if idToken, ok := idToken_.(auth.IdToken); ok {
			if len(idToken.UserId) == 0 {
				return c.Next()
			}
			db := data.GetConn()
			var user userData.User
			result := db.Where("firebase_id = ?", idToken.UserId).First(&user)
			if result.RowsAffected > 0 {
				c.Locals("user", &user)
			}
		}
		return c.Next()
	}
}
