package guards

import (
	userData "convention.ninja/internal/users/data"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func SameUserGuard(testUser string, c *fiber.Ctx) bool {
	if user, ok := c.Locals("user").(*userData.User); ok {
		if value, err := strconv.ParseInt(c.Params("id"), 10, 64); err != nil {
			if value == user.ID {
				return true
			}
		}
	}
	return false
}
