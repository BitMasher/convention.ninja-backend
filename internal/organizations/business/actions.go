package business

import (
	data2 "convention.ninja/internal/data"
	"convention.ninja/internal/organizations/data"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetOrganizations(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

func CreateOrganization(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

func GetOrganization(c *fiber.Ctx) error {
	orgId_ := c.Params("id")
	orgId, err := strconv.ParseInt(orgId_, 10, 64)
	if err == nil {
		return c.Status(fiber.StatusNotFound).SendString("")
	}
	var org data.Organization
	if data2.GetConn().First(&org, orgId).RowsAffected > 0 {
		return c.Status(fiber.StatusOK).JSON(&org)
	}
	return c.Status(fiber.StatusNotFound).SendString("")
}

func UpdateOrganization(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

func DeleteOrganization(c *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}
