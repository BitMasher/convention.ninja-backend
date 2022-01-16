package business

import (
	"convention.ninja/internal/common"
	data2 "convention.ninja/internal/data"
	data3 "convention.ninja/internal/inventory/data"
	"convention.ninja/internal/snowflake"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetManufacturers(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var manufacturers []data3.Manufacturer
	data2.GetConn().Where(&data3.Manufacturer{OrganizationId: org.ID}).Find(&manufacturers)
	return c.Status(fiber.StatusOK).JSON(&manufacturers)
}

type CreateManufacturerRequest struct {
	Name string `json:"name"`
}

func CreateManufacturer(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var req CreateManufacturerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	count := int64(0)
	data2.GetConn().Where(&data3.Manufacturer{Name: req.Name, OrganizationId: org.ID}).Count(&count)
	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	mfg := data3.Manufacturer{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:           req.Name,
		OrganizationId: org.ID,
	}
	data2.GetConn().Create(&mfg)
	return c.Status(fiber.StatusOK).JSON(&mfg)
}

func GetManufacturer(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	mfgId_ := c.Params("mfgId", "")
	if mfgId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	mfgId, err := strconv.ParseInt(mfgId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var mfg data3.Manufacturer
	if data2.GetConn().Where(&data3.Manufacturer{
		OrganizationId: org.ID,
	}).First(&mfg, mfgId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&mfg)
}

type UpdateManufacturerRequest struct {
	Name string `json:"name"`
}

func UpdateManufacturer(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var req UpdateManufacturerRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	mfgId_ := c.Params("mfgId", "")
	if mfgId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	mfgId, err := strconv.ParseInt(mfgId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var mfg data3.Manufacturer
	if data2.GetConn().Where(&data3.Manufacturer{
		OrganizationId: org.ID,
	}).First(&mfg, mfgId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	count := int64(0)
	data2.GetConn().Where(&data3.Manufacturer{Name: req.Name, OrganizationId: org.ID}).Count(&count)
	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}
	mfg.Name = req.Name
	data2.GetConn().Save(&mfg)
	return c.Status(fiber.StatusOK).JSON(&mfg)
}

func DeleteManufacturer(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	mfgId_ := c.Params("mfgId", "")
	if mfgId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	mfgId, err := strconv.ParseInt(mfgId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var mfg data3.Manufacturer
	if data2.GetConn().Where(&data3.Manufacturer{
		OrganizationId: org.ID,
	}).First(&mfg, mfgId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data2.GetConn().Delete(&mfg)
	return c.SendStatus(fiber.StatusOK)
}
