package business

import (
	"convention.ninja/internal/common"
	data2 "convention.ninja/internal/data"
	"convention.ninja/internal/inventory/data"
	"convention.ninja/internal/snowflake"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func GetAssets(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	var assets []data.Asset
	data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).
		Preload("AssetTags").
		Joins("Model").
		Joins("Model.Category").
		Joins("Model.Manufacturer").Find(&assets)
	return c.Status(fiber.StatusOK).JSON(&assets)
}

type CreateAssetRequest struct {
	SerialNumber string   `json:"serialNumber"`
	ModelId      int64    `json:"modelId"`
	AssetTags    []string `json:"assetTags"`
}

func CreateAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	var req CreateAssetRequest
	err := c.BodyParser(&req)
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.ModelId == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	asset := &data.Asset{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		ModelId:        req.ModelId,
		SerialNumber:   req.SerialNumber,
		OrganizationId: org.ID,
	}
	db := data2.GetConn()
	var count int64
	if db.Where(&data.Model{
		SnowflakeModel: data2.SnowflakeModel{ID: req.ModelId},
		OrganizationId: org.ID,
	}).Count(&count); count == 0 {
		return c.SendStatus(fiber.StatusBadRequest)
	}
	if req.AssetTags != nil && len(req.AssetTags) > 0 {
		if db.Model(&data.AssetTag{}).Where("organization_id = ? AND tag_id in ?", org.ID, req.AssetTags).Count(&count); count > 0 {
			return c.SendStatus(fiber.StatusConflict)
		}
	}
	db.Create(&asset)
	asset.AssetTags = make([]data.AssetTag, 0)
	for _, tag := range req.AssetTags {
		assetTag := data.AssetTag{
			SnowflakeModel: data2.SnowflakeModel{ID: snowflake.GetNode().Generate().Int64()},
			TagId:          tag,
			AssetId:        asset.ID,
			OrganizationId: org.ID,
		}
		db.Create(&assetTag)
		asset.AssetTags = append(asset.AssetTags, assetTag)
	}
	return c.Status(fiber.StatusOK).JSON(&asset)
}

func GetAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).
		Preload("AssetTags").
		Joins("Model").
		Joins("Model.Category").
		Joins("Model.Manufacturer").First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&asset)
}

type UpdateAssetRequest struct {
	SerialNumber string `json:"serialNumber"`
	ModelId      int64
}

func UpdateAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
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
	db := data2.GetConn()
	var count int64
	if req.ModelId > 0 {
		if db.Where(&data.Model{
			SnowflakeModel: data2.SnowflakeModel{ID: req.ModelId},
			OrganizationId: org.ID,
		}).Count(&count); count == 0 {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if req.SerialNumber != "" {
		asset.SerialNumber = req.SerialNumber
	}
	if req.ModelId != 0 {
		asset.ModelId = req.ModelId
	}

	db.Save(&asset)

	return c.Status(fiber.StatusOK).JSON(&asset)
}

func GetAssetBarcodes(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).Preload("AssetTags").First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&asset.AssetTags)
}

type CreateAssetBarcodeRequest struct {
	TagId string `json:"tagId"`
}

func CreateAssetBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
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
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var count int64
	if data2.GetConn().Where(&data.AssetTag{
		TagId:          req.TagId,
		OrganizationId: org.ID,
	}).Count(&count); count > 0 {
		return c.SendStatus(fiber.StatusConflict)
	}
	tag := data.AssetTag{
		SnowflakeModel: data2.SnowflakeModel{
			ID: snowflake.GetNode().Generate().Int64(),
		},
		TagId:          req.TagId,
		AssetId:        asset.ID,
		OrganizationId: org.ID,
	}
	data2.GetConn().Create(&tag)
	return c.Status(fiber.StatusOK).JSON(&tag)
}

func DeleteAssetBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
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
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var assetTag data.AssetTag
	if data2.GetConn().Where(&data.AssetTag{
		AssetId:        asset.ID,
		OrganizationId: org.ID,
	}).First(&assetTag, barcodeId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data2.GetConn().Delete(&assetTag)
	return c.SendStatus(fiber.StatusOK)
}

func DeleteAsset(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	assetId_ := c.Params("assetId", "")
	if assetId_ == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	assetId, err := strconv.ParseInt(assetId_, 10, 64)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var asset data.Asset
	if data2.GetConn().Where(&data.Asset{
		OrganizationId: org.ID,
	}).Preload("AssetTags").First(&asset, assetId).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	data2.GetConn().Where(data.AssetTag{
		AssetId:        asset.ID,
		OrganizationId: org.ID,
	}).Delete(data.AssetTag{})
	data2.GetConn().Delete(&asset)
	return c.SendStatus(fiber.StatusOK)
}

func GetAssetByBarcode(c *fiber.Ctx) error {
	org, auth := common.GetOrgAndAuthorize(c)
	if org == nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	if auth == false {
		return c.SendStatus(fiber.StatusUnauthorized)
	}
	tagId := c.Params("barcode", "")
	if tagId == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	var tag data.AssetTag
	if data2.GetConn().Joins("Asset").Where(&data.AssetTag{
		TagId:          tagId,
		OrganizationId: org.ID,
	}).First(&tag).RowsAffected == 0 {
		return c.SendStatus(fiber.StatusNotFound)
	}
	return c.Status(fiber.StatusOK).JSON(&tag.Asset)
}
