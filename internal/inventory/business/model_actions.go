package business

import (
	"convention.ninja/internal/common"
	"convention.ninja/internal/inventory/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetModels(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	models, err := data.GetModelsExpandedByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetModels: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&models)
}

type CreateModelRequest struct {
	Name           string `json:"name"`
	ManufacturerId int64  `json:"manufacturerId,string"`
	CategoryId     int64  `json:"categoryId,string"`
}

func CreateModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var req CreateModelRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	exists, err := data.ModelExistsInOrg(org.ID, req.Name, req.ManufacturerId)
	if err != nil {
		fmt.Printf("got error in CreateModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}

	exists, err = data.ManufacturerExistsById(req.ManufacturerId, org.ID)
	if err != nil {
		fmt.Printf("got error in CreateModel: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	exists, err = data.CategoryExistsById(req.CategoryId, org.ID)
	if err != nil {
		fmt.Printf("got error in CreateModel: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	model := data.Model{
		Name:           req.Name,
		ManufacturerId: req.ManufacturerId,
		CategoryId:     req.CategoryId,
		OrganizationId: org.ID,
	}

	err = data.CreateModel(&model)
	if err != nil {
		fmt.Printf("got error in CreateModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&model)
}

func GetModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId_ := c.Params("modelId", "")
	if modelId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId, err := strconv.ParseInt(modelId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	model, err := data.GetModelExpandedById(modelId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if model == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&model)
}

type UpdateModelRequest struct {
	Name           string `json:"name"`
	ManufacturerId int64  `json:"manufacturerId,string"`
	CategoryId     int64  `json:"categoryId,string"`
}

func UpdateModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var req UpdateModelRequest
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

	model, err := data.GetModelById(modelId, org.ID)
	if err != nil {
		fmt.Printf("got error in UpdateModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if model == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	if len(req.Name) > 0 {
		tmpMfg := model.ManufacturerId
		if req.ManufacturerId > 0 {
			tmpMfg = req.ManufacturerId
		}
		exists, err := data.ModelExistsInOrg(org.ID, req.Name, tmpMfg)
		if err != nil {
			fmt.Printf("got error in UpdateModel: %s\n", err) // TODO Implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if exists {
			return c.SendStatus(fiber.StatusConflict)
		}
		model.Name = req.Name
	}

	if req.ManufacturerId > 0 {
		exists, err := data.ManufacturerExistsById(req.ManufacturerId, org.ID)
		if err != nil {
			fmt.Printf("got error in UpdateModel: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !exists {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		tmpName := model.Name
		if len(req.Name) > 0 {
			tmpName = req.Name
		}
		exists, err = data.ModelExistsInOrg(org.ID, tmpName, req.ManufacturerId)
		model.ManufacturerId = req.ManufacturerId
	}
	if req.CategoryId > 0 {
		exists, err := data.CategoryExistsById(req.CategoryId, org.ID)
		if err != nil {
			fmt.Printf("got error in UpdateModel: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !exists {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		model.CategoryId = req.CategoryId
	}
	err = data.UpdateModel(model)
	if err != nil {
		fmt.Printf("got error in UpdateModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&model)
}

func DeleteModel(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	modelId_ := c.Params("modelId", "")
	if modelId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	modelId, err := strconv.ParseInt(modelId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	model, err := data.GetModelById(modelId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if model == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteModel(model)
	if err != nil {
		fmt.Printf("got error in DeleteModel: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.SendStatus(fiber.StatusOK)
}
