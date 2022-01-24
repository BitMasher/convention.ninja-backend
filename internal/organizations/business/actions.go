package business

import (
	"convention.ninja/internal/common"
	"convention.ninja/internal/organizations/data"
	userData "convention.ninja/internal/users/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
)

func GetOrganizations(c *fiber.Ctx) error {
	// TODO this should have search filters instead of just returning _my_ organizations
	user := c.Locals("user").(*userData.User)
	organizations, err := data.GetOrganizationsByOwner(user.ID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&organizations)
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
	orgExists, err := data.OrganizationNameExists(req.Name)
	if err != nil {
		fmt.Printf("got error in CreateOrganization: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if orgExists {
		return c.SendStatus(fiber.StatusConflict)
	}
	org := data.Organization{
		Name:    req.Name,
		OwnerId: c.Locals("user").(*userData.User).ID,
	}

	if err = data.CreateOrganization(&org); err != nil {
		fmt.Printf("got error in CreateOrganization: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&org)
}

func GetOrganization(c *fiber.Ctx) error {
	orgId_ := c.Params("orgId")
	orgId, err := strconv.ParseInt(orgId_, 10, 64)
	if err == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	org, err := data.GetOrganizationById(orgId)
	if err != nil {
		fmt.Printf("got error in GetOrganization: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if org != nil {
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
		exists, err := data.OrganizationNameExists(req.Name)
		if err != nil {
			fmt.Printf("got error in UpdateOrganizatin: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if exists {
			return c.SendStatus(fiber.StatusConflict)
		}
	}
	org.Name = req.Name
	if err = data.UpdateOrganization(org); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&org)
}

func DeleteOrganization(c *fiber.Ctx) error {
	// TODO this needs to do a lot more clean up
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if err := data.DeleteOrganization(org); err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
