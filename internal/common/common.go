package common

import (
	"convention.ninja/internal/auth/guards"
	data2 "convention.ninja/internal/data"
	"convention.ninja/internal/organizations/data"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetOrgAndAuthorize(c *fiber.Ctx) (*data.Organization, bool) {
	var org data.Organization
	orgId_ := c.Params("orgId", "")
	if orgId_ == "" {
		return nil, true
	}
	orgId, err := strconv.ParseInt(orgId_, 10, 64)
	if err != nil {
		return nil, true
	}
	if data2.GetConn().First(&org, orgId).RowsAffected == 0 {
		return nil, true
	}
	if !guards.IsAuthorizedToOrg(&org, c) {
		return &org, false
	}
	return &org, true
}
