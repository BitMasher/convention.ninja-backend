package business

import (
	"convention.ninja/internal/common"
	"convention.ninja/internal/inventory/data"
	userData "convention.ninja/internal/users/data"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetManifests(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifests, err := data.GetOpenManifestsByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetManifests: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&manifests)
}

func CreateManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	manifest := data.Manifest{
		CreatorId:      c.Locals("user").(*userData.User).ID,
		OrganizationId: org.ID,
	}
	err := data.CreateManifest(&manifest)
	if err != nil {
		fmt.Printf("got error in CreateManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&manifest)
}

func GetManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&manifest)
}

type UpdateManifestRequest struct {
	RoomId                   string             `json:"roomId"`
	ResponsibleExternalParty data.ExternalParty `json:"responsibleExternalParty"`
}

func UpdateManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req UpdateManifestRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in UpdateManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	changed := false
	if manifest.RoomId.String != req.RoomId {
		manifest.RoomId = sql.NullString{
			String: req.RoomId,
			Valid:  len(req.RoomId) > 0,
		}
		changed = true
	}
	if manifest.ResponsibleExternalParty.Name != req.ResponsibleExternalParty.Name || manifest.ResponsibleExternalParty.Extra != req.ResponsibleExternalParty.Extra {
		manifest.ResponsibleExternalParty.Name = req.ResponsibleExternalParty.Name
		manifest.ResponsibleExternalParty.Extra = req.ResponsibleExternalParty.Extra
		manifest.ResponsibleExternalParty.Valid = len(req.ResponsibleExternalParty.Name) > 0 || len(req.ResponsibleExternalParty.Extra) > 0
		changed = true
	}
	if changed {
		err = data.UpdateManifest(manifest)
		if err != nil {
			fmt.Printf("got error in UpdateManifest: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}
	return c.Status(fiber.StatusOK).JSON(&manifest)
}

func ShipManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in ShipManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if manifest.ShipDate.Valid {
		return c.Status(fiber.StatusOK).JSON(&manifest)
	}
	err = data.ShipManifest(manifest)
	if err != nil {
		fmt.Printf("got error in ShipManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&manifest)
}

func UnshipManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in UnshipManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if !manifest.ShipDate.Valid {
		return c.Status(fiber.StatusOK).JSON(&manifest)
	}
	err = data.UnshipManifest(manifest)
	if err != nil {
		fmt.Printf("got error in UnshipManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&manifest)
}

func DeleteManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteManifest(manifest)
	if err != nil {
		fmt.Printf("got error in DeleteManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetManifestAssets(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetManifestAssets: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	entries, err := data.GetManifestEntriesExpanded(manifest.ID, org.ID)
	if err != nil {
		fmt.Printf("got error in GetManifestAssets: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&entries)
}

type AddManifestEntryRequest struct {
	AssetId int64  `json:"assetId,string"`
	TagId   string `json:"tagId"`
}

func AddAssetToManifest(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in AddAssetToManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req AddManifestEntryRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.AssetId == 0 && len(req.TagId) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.AssetId != 0 && len(req.TagId) > 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	manifestEntry := data.ManifestEntry{
		ManifestId:     manifest.ID,
		OrganizationId: org.ID,
	}
	if req.AssetId != 0 {
		asset, err := data.GetAssetById(req.AssetId, org.ID)
		if err != nil {
			fmt.Printf("got error in AddAssetToManifest: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if asset == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		manifestEntry.AssetId = asset.ID
		manifestEntry.Asset = *asset

	}
	if len(req.TagId) > 0 {
		asset, err := data.GetAssetByTag(req.TagId, org.ID)
		if err != nil {
			fmt.Printf("got error in AddAssetToManifest: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if asset == nil {
			return c.SendStatus(fiber.StatusBadRequest)
		}
		manifestEntry.AssetId = asset.ID
		manifestEntry.Asset = *asset
	}
	exists, err := data.AssetExistsInManifest(org.ID, manifest.ID, manifestEntry.AssetId)
	if err != nil {
		fmt.Printf("got error in AddAssetToManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}
	err = data.AddEntryToManifest(&manifestEntry)
	if err != nil {
		fmt.Printf("got error in AddAssetToManifest: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&manifestEntry)
}

func DeleteManifestEntry(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId_ := c.Params("manifestId", "")
	if manifestId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestId, err := strconv.ParseInt(manifestId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifest, err := data.GetManifestById(manifestId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteManifestEntry: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifest == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	entryId_ := c.Params("entryId", "")
	if entryId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	entryId, err := strconv.ParseInt(entryId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	manifestEntry, err := data.GetManifestEntryById(manifest.ID, entryId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteManifestEntry: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if manifestEntry == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteManifestEntry(manifestEntry)
	if err != nil {
		fmt.Printf("got error in DeleteManifestEntry: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}
