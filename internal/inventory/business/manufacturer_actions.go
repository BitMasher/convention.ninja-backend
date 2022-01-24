package business

import (
	"convention.ninja/internal/common"
	data3 "convention.ninja/internal/inventory/data"
	"fmt"
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
	manufacturers, err := data3.GetManufacturersByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetManufacturers: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
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
	exists, err := data3.ManufacturerExistsInOrg(org.ID, req.Name)
	if err != nil {
		fmt.Printf("got error in CreateManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}

	mfg := data3.Manufacturer{
		Name:           req.Name,
		OrganizationId: org.ID,
	}
	err = data3.CreateManufacturer(&mfg)
	if err != nil {
		fmt.Printf("got error in CreateManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
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
	mfg, err := data3.GetManufacturerById(mfgId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if mfg == nil {
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
	mfg, err := data3.GetManufacturerById(mfgId, org.ID)
	if err != nil {
		fmt.Printf("got error in UpdateManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if mfg == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	exists, err := data3.ManufacturerExistsInOrg(org.ID, req.Name)
	if err != nil {
		fmt.Printf("got error in UpdateManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}
	mfg.Name = req.Name
	err = data3.UpdateManufacturer(mfg)
	if err != nil {
		fmt.Printf("got error in UpdateManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
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
	mfg, err := data3.GetManufacturerById(mfgId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteManufacturer: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if mfg == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data3.DeleteManufacturer(mfg)
	if err != nil {
		fmt.Printf("got error in DeleteManufacturer: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
