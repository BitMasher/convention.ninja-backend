package business

import (
	"convention.ninja/internal/common"
	"convention.ninja/internal/inventory/data"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAssets(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assets, err := data.GetAssetsExpandedByOrganization(org.ID)
	if err != nil {
		fmt.Printf("got error in GetAssets: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&assets)
}

type CreateAssetRequest struct {
	SerialNumber string   `json:"serialNumber"`
	ModelId      int64    `json:"modelId,string"`
	AssetTags    []string `json:"assetTags"`
}

func CreateAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var req CreateAssetRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.ModelId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	asset := data.Asset{
		ModelId:        req.ModelId,
		SerialNumber:   req.SerialNumber,
		OrganizationId: org.ID,
	}
	exists, err := data.ModelExistsById(req.ModelId, org.ID)
	if err != nil {
		fmt.Printf("got error in CreateAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.AssetTags != nil && len(req.AssetTags) > 0 {
		exists, err = data.AssetTagsExistInOrg(org.ID, req.AssetTags)
		if err != nil {
			fmt.Printf("got error in CreateAsset: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if exists {
			return c.SendStatus(fiber.StatusConflict)
		}
	}
	err = data.CreateAsset(&asset)
	if err != nil {
		fmt.Printf("got error in CreateAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	asset.AssetTags = make([]data.AssetTag, 0)
	for _, tag := range req.AssetTags {
		assetTag := data.AssetTag{
			TagId:          tag,
			AssetId:        asset.ID,
			OrganizationId: org.ID,
		}
		asset.AssetTags = append(asset.AssetTags, assetTag)
	}
	err = data.BulkCreateAssetTag(asset.AssetTags)
	if err != nil {
		fmt.Printf("got error in CreateAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&asset)
}

func GetAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	asset, err := data.GetAssetExpandedById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if asset == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&asset)
}

type UpdateAssetRequest struct {
	SerialNumber string `json:"serialNumber"`
	ModelId      int64  `json:"modelId,string"`
	RoomId       string `json:"roomId"`
}

func UpdateAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req UpdateAssetRequest
	err = c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.SerialNumber == "" && req.ModelId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.ModelId > 0 {
		exists, err := data.ModelExistsById(req.ModelId, org.ID)
		if err != nil {
			fmt.Printf("got error in UpdateAsset: %s\n", err) // TODO implement logging system
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		if !exists {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	asset, err := data.GetAssetById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in UpdateAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if asset == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if req.SerialNumber != "" {
		asset.SerialNumber = req.SerialNumber
	}
	if req.ModelId != 0 {
		asset.ModelId = req.ModelId
	}
	if req.RoomId != "" {
		asset.RoomId = req.RoomId
	}
	err = data.UpdateAsset(asset)
	if err != nil {
		fmt.Printf("got error in UpdateAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(&asset)
}

func GetAssetBarcodes(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	exists, err := data.AssetExistsById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetAssetBarcodes: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	tags, err := data.GetAssetTagsByAssetId(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetAssetBarcodes: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&tags)
}

type CreateAssetBarcodeRequest struct {
	TagId string `json:"tagId"`
}

func CreateAssetBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var req CreateAssetBarcodeRequest
	err = c.BodyParser(&req)
	if err != nil || len(req.TagId) == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	exists, err := data.AssetExistsById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in CreateAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	exists, err = data.AssetTagExistInOrg(org.ID, req.TagId)
	if err != nil {
		fmt.Printf("got error in CreateAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if exists {
		return c.SendStatus(fiber.StatusConflict)
	}
	tag := data.AssetTag{
		TagId:          req.TagId,
		AssetId:        assetId,
		OrganizationId: org.ID,
	}
	err = data.CreateAssetTag(&tag)
	if err != nil {
		fmt.Printf("got error in CreateAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Status(fiber.StatusOK).JSON(&tag)
}

func DeleteAssetBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	barcodeId_ := c.Params("barcodeId", "")
	if barcodeId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	barcodeId, err := strconv.ParseInt(barcodeId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	exists, err := data.AssetExistsById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if !exists {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetTag, err := data.GetAssetTagById(barcodeId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if assetTag == nil || assetTag.AssetId != assetId {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteAssetTag(assetTag)
	if err != nil {
		fmt.Printf("got error in DeleteAssetBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	asset, err := data.GetAssetById(assetId, org.ID)
	if err != nil {
		fmt.Printf("got error in DeleteAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if asset == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	err = data.DeleteAsset(asset)
	if err != nil {
		fmt.Printf("got error in DeleteAsset: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.SendStatus(fiber.StatusOK)
}

func GetAssetByBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	tagId := c.Params("barcode", "")
	if tagId == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	asset, err := data.GetAssetExpandedByTag(tagId, org.ID)
	if err != nil {
		fmt.Printf("got error in GetAssetByBarcode: %s\n", err) // TODO implement logging system
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	if asset == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&asset)
}
