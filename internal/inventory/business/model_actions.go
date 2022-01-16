package business

import (
	"convention.ninja/internal/common"
	data2 "convention.ninja/internal/data"
	"convention.ninja/internal/inventory/data"
	"convention.ninja/internal/snowflake"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetModels(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var models []data.Model
	data2.GetConn().Where(&data.Model{
		OrganizationId: org.ID,
	}).Joins("Category").Joins("Manufacturer").Find(&models)
	return c.Status(fiber.StatusOK).JSON(&models)
}

type CreateModelRequest struct {
	Name           string `json:"name"`
	ManufacturerId int64  `json:"manufacturerId"`
	CategoryId     int64  `json:"categoryId"`
}

func CreateModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateModelRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	var count int64
	if data2.GetConn().Where(&data.Model{OrganizationId: org.ID, Name: req.Name}).Count(&count); count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}

	if data2.GetConn().Where(&data.Manufacturer{OrganizationId: org.ID, SnowflakeModel: data2.SnowflakeModel{ID: req.ManufacturerId}}).Count(&count); count == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if data2.GetConn().Where(&data.Category{OrganizationId: org.ID, SnowflakeModel: data2.SnowflakeModel{ID: req.ManufacturerId}}).Count(&count); count == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	model := &data.Model{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		Name:           req.Name,
		ManufacturerId: req.ManufacturerId,
		CategoryId:     req.CategoryId,
		OrganizationId: org.ID,
	}

	data2.GetConn().Create(&model)

	return c.Status(fiber.StatusOK).JSON(&model)
}

func GetModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	modelId_ := c.Params("modelId", "")
	if modelId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId, err := strconv.ParseInt(modelId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var model data.Model
	if data2.GetConn().Where(&data.Model{OrganizationId: org.ID}).Joins("Manufacturer").Joins("Category").First(&model, modelId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&model)
}

type UpdateModelRequest struct {
	Name           string `json:"name"`
	ManufacturerId int64  `json:"manufacturerId"`
	CategoryId     int64  `json:"categoryId"`
}

func UpdateModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateModelRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if len(req.Name) == 0 && req.CategoryId == 0 && req.ManufacturerId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	modelId_ := c.Params("modelId", "")
	if modelId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId, err := strconv.ParseInt(modelId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var model data.Model
	if data2.GetConn().Where(&data.Model{
		OrganizationId: org.ID}).First(&model, modelId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var count int64
	if len(req.Name) > 0 {
		if data2.GetConn().Where(&data.Model{
			OrganizationId: org.ID, Name: req.Name}).Count(&count); count > 0 {
			return c.SendStatus(fiber.StatusConflict)
		}
		model.Name = req.Name
	}

	if req.ManufacturerId > 0 {
		if data2.GetConn().Where(&data.Manufacturer{
			OrganizationId: org.ID, SnowflakeModel: data2.SnowflakeModel{
				ID: req.ManufacturerId}}).Count(&count); count == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		model.ManufacturerId = req.ManufacturerId
	}
	if req.CategoryId > 0 {
		if data2.GetConn().Where(&data.Category{
			OrganizationId: org.ID, SnowflakeModel: data2.SnowflakeModel{
				ID: req.ManufacturerId}}).Count(&count); count == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		model.CategoryId = req.CategoryId
	}
	data2.GetConn().Save(&model)

	return c.Status(fiber.StatusOK).JSON(&model)
}

func DeleteModeL(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	modelId_ := c.Params("modelId", "")
	if modelId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId, err := strconv.ParseInt(modelId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var model data.Model
	if data2.GetConn().Where(&data.Model{
		OrganizationId: org.ID}).First(&model, modelId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data2.GetConn().Delete(&model)

	return c.SendStatus(fiber.StatusOK)
}
