package guards

import (
	"convention.ninja/internal/organizations/data"
	userData "convention.ninja/internal/users/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func SameUserGuard(testUser string, c *fiber.Ctx) bool {
	if user, ok := c.Locals("user").(*userData.User); ok {
		value, err := strconv.ParseInt(testUser, 10, 64)
		if err != nil {
			fmt.Printf("got error in SameUserGuard: %s\n", err) // TODO implement logging system
			return false
		}
		return value == user.ID
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
