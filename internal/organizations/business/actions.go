package business

import (
	"convention.ninja/internal/common"
	data2 "convention.ninja/internal/data"
	"convention.ninja/internal/organizations/data"
	"convention.ninja/internal/snowflake"
	userData "convention.ninja/internal/users/data"
	"errors"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func GetOrganizations(_ *fiber.Ctx) error {
	return errors.New("not implemented") // TODO
}

type CreateOrganizationRequest struct {
	Name string `json:"name"`
}

func CreateOrganization(c *fiber.Ctx) error {
	var req CreateOrganizationRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	var count int64
	if data2.GetConn().Where(&data.Organization{NormalizedName: strings.ToLower(req.Name)}).Count(&count); count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}
	org := data.Organization{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:           req.Name,
		NormalizedName: strings.ToLower(req.Name),
		OwnerId:        c.Locals("user").(*userData.User).ID,
	}

	data2.GetConn().Create(&org)

	return c.Status(fiber.StatusOK).JSON(&org)
}

func GetOrganization(c *fiber.Ctx) error {
	orgId_ := c.Params("id")
	orgId, err := strconv.ParseInt(orgId_, 10, 64)
	if err == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var org data.Organization
	if data2.GetConn().First(&org, orgId).RowsAffected > 0 {
		return c.Status(fiber.StatusOK).JSON(&org)
	}
	return c.SendStatus(fiber.StatusNotFound)
}

type UpdateOrganizationRequest struct {
	Name string `json:"name"`
}

func UpdateOrganization(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req UpdateOrganizationRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if org.NormalizedName != strings.ToLower(req.Name) {
		var count int64
		if data2.GetConn().Where(&data.Organization{
			NormalizedName: strings.ToLower(req.Name)}).Count(&count); count > 0 {
			return c.SendStatus(fiber.StatusConflict)
		}
	}
	org.Name = req.Name
	org.NormalizedName = strings.ToLower(req.Name)
	data2.GetConn().Save(&org)
	return c.Status(fiber.StatusOK).JSON(&org)
}

func DeleteOrganization(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	data2.GetConn().Delete(&org)
	return c.SendStatus(fiber.StatusOK)
}
