package business

import (
	"convention.ninja/internal/common"
	data2 "convention.ninja/internal/data"
	data3 "convention.ninja/internal/inventory/data"
	"convention.ninja/internal/snowflake"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetCategories(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var categories []data3.Category
	data2.GetConn().Where(&data3.Category{OrganizationId: org.ID}).Find(&categories)
	return c.Status(fiber.StatusOK).JSON(&categories)
}

type CreateCategoryRequest struct {
	Name string
}

func CreateCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var req CreateCategoryRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	count := int64(0)
	data2.GetConn().Where(&data3.Category{Name: req.Name, OrganizationId: org.ID}).Count(&count)
	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	cat := data3.Category{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:           req.Name,
		OrganizationId: org.ID,
	}
	data2.GetConn().Create(&cat)
	return c.Status(fiber.StatusOK).JSON(&cat)
}

func GetCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	catId_ := c.Params("categoryId", "")
	if catId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	catId, err := strconv.ParseInt(catId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var mfg data3.Category
	if data2.GetConn().Where(&data3.Category{
		OrganizationId: org.ID,
	}).First(&mfg, catId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&mfg)
}

type UpdateCategoryRequest struct {
	Name string
}

func UpdateCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var req UpdateCategoryRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	catId_ := c.Params("categoryId", "")
	if catId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	catId, err := strconv.ParseInt(catId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var cat data3.Category
	if data2.GetConn().Where(&data3.Category{
		OrganizationId: org.ID,
	}).First(&cat, catId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	count := int64(0)
	data2.GetConn().Where(&data3.Category{Name: req.Name, OrganizationId: org.ID}).Count(&count)
	if count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}
	cat.Name = req.Name
	data2.GetConn().Save(&cat)
	return c.Status(fiber.StatusOK).JSON(&cat)
}

func DeleteCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	catId_ := c.Params("categoryId", "")
	if catId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	catId, err := strconv.ParseInt(catId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var cat data3.Category
	if data2.GetConn().Where(&data3.Category{
		OrganizationId: org.ID,
	}).First(&cat, catId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data2.GetConn().Delete(&cat)
	return c.SendStatus(fiber.StatusOK)
}
