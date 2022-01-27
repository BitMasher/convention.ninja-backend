package common

import (
	"convention.ninja/internal/auth/guards"
	"convention.ninja/internal/organizations/data"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetOrgAndAuthorize(c *fiber.Ctx) (*data.Organization, bool) {
	orgId_ := c.Params("orgId", "")
	if orgId_ == "" {
		return nil, false
	}
	orgId, err := strconv.ParseInt(orgId_, 10, 64)
	if err != nil {
		return nil, false
	}
	org, err := data.GetOrganizationById(orgId)
	if org == nil {
		return nil, false
	}
	if !guards.IsAuthorizedToOrg(org, c) {
		return org, false
	}
	return org, true
}
