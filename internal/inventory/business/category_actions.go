package business

import (
	"convention.ninja/internal/common"
	data3 "convention.ninja/internal/inventory/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetCategories(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	categories, err := data3.GetCategoriesByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetCategories: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&categories)
}

type CreateCategoryRequest struct {
	Name string `json:"name"`
}

func CreateCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req CreateCategoryRequest
	err := c.BodyParser(&req)
	if err != nil || len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	exists, err := data3.CategoryExistsInOrg(org.ID, req.Name)
	if err != nil {
		fmt.Printf("got error in CreateCategory: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}

	cat := data3.Category{
		Name:           req.Name,
		OrganizationId: org.ID,
	}
	err = data3.CreateCategory(&cat)
	if err != nil {
		fmt.Printf("got error in CreateCategory: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&cat)
}

func GetCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	catId_ := c.Params("categoryId", "")
	if catId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	catId, err := strconv.ParseInt(catId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	cat, err := data3.GetCategoryById(catId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetCategory: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if cat == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&cat)
}

type UpdateCategoryRequest struct {
	Name string `json:"name"`
}

func UpdateCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req UpdateCategoryRequest
	err := c.BodyParser(&req)
	if err != nil || len(req.Name) == 0 {
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
	cat, err := data3.GetCategoryById(catId, org.ID)
	if err != nil {
		fmt.Printf("got error in UpdateCategory: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if cat == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	exists, err := data3.CategoryExistsInOrg(org.ID, req.Name)
	if err != nil {
		fmt.Printf("got error in UpdateCategory: %s\n", err) // TODO implement logging system`
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}
	cat.Name = req.Name
	err = data3.UpdateCategory(cat)
	if err != nil {
		fmt.Printf("got error in UpdateCategory: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&cat)
}

func DeleteCategory(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	catId_ := c.Params("categoryId", "")
	if catId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	catId, err := strconv.ParseInt(catId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	cat, err := data3.GetCategoryById(catId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteCategory: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if cat == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data3.DeleteCategory(cat)
	if err != nil {
		fmt.Printf("got error in UpdateCategory: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
