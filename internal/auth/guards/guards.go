package guards

import (
	"convention.ninja/internal/organizations/data"
	userData "convention.ninja/internal/users/data"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func SameUserGuard(testUser string, c *fiber.Ctx) bool {
	if user, ok := c.Locals("user").(*userData.User); ok {
		if value, err := strconv.ParseInt(testUser, 10, 64); err != nil {
			if value == user.ID {
				return true
			}
		}
	}
	return false
}

func IsAuthorizedToOrg(org *data.Organization, c *fiber.Ctx) bool {
	// TODO implement RBAC for organizations
	if org.OwnerId == c.Locals("user").(*userData.User).ID {
		return true
	}
	return false
}
