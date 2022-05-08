package business

import (
	"convention.ninja/internal/common"
	"convention.ninja/internal/organizations/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetVenues(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venues, err := data.GetVenuesByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenues: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&venues)
}

type CreateVenueRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func CreateVenue(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req CreateVenueRequest
	err := c.BodyParser(&req)
	if err != nil || len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	venue := data.Venue{
		Name:           req.Name,
		OrganizationId: org.ID,
		Address:        req.Address,
		Rooms:          make([]data.VenueRoom, 0),
	}
	err = data.CreateVenue(&venue)
	if err != nil {
		fmt.Printf("got error in CreateVenue: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&venue)
}

func GetVenue(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId_ := c.Params("venueId", "")
	if venueId_ == "" {
		c.SendStatus(fiber.StatusNotFound)
	}
	venueId, err := strconv.ParseInt(venueId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venue, err := data.GetVenueById(venueId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenue: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&venue)
}

type UpdateVenueRequest struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func UpdateVenue(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req UpdateVenueRequest
	err := c.BodyParser(&req)
	if err != nil || (len(req.Name) == 0 && len(req.Address) == 0) {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	venueId_ := c.Params("venueId", "")
	if venueId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId, err := strconv.ParseInt(venueId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venue, err := data.GetVenueById(venueId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenue: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if venue == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if len(req.Name) > 0 {
		venue.Name = req.Name
	}
	if len(req.Address) > 0 {
		venue.Address = req.Address
	}
	err = data.UpdateVenue(venue)
	if err != nil {
		fmt.Printf("got error in GetVenue: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&venue)
}

func DeleteVenue(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId_ := c.Params("venueId", "")
	if venueId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId, err := strconv.ParseInt(venueId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venue, err := data.GetVenueById(venueId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenue: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if venue == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteVenue(venue)
	if err != nil {
		fmt.Printf("got error in DeleteVenue: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetVenueRooms(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId_ := c.Params("venueId", "")
	if venueId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId, err := strconv.ParseInt(venueId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	rooms, err := data.GetVenueRooms(venueId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenueRooms: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&rooms)
}

type CreateVenueRoomRequest struct {
	Name string `json:"name"`
}

func CreateVenueRoom(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req CreateVenueRoomRequest
	err := c.BodyParser(&req)
	if err != nil || len(req.Name) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	venueId_ := c.Params("venueId", "")
	if venueId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueId, err := strconv.ParseInt(venueId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venue, err := data.GetVenueById(venueId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetVenue: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if venue == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	room := data.VenueRoom{
		Name:           req.Name,
		VenueId:        venue.ID,
		OrganizationId: org.ID,
	}
	err = data.CreateVenueRoom(&room)
	if err != nil {
		fmt.Printf("got error in CreateVenueRoom: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&room)
}

func DeleteVenueRoom(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	roomId_ := c.Params("roomId", "")
	if roomId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	roomId, err := strconv.ParseInt(roomId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	venueRoom, err := data.GetVenueRoomById(roomId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteVenueRoom: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if venueRoom == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteVenueRoom(venueRoom)
	if err != nil {
		fmt.Printf("got error in DeleteVenueRoom: %s\n", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
