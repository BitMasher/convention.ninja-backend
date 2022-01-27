package inventory

import (
	"convention.ninja/internal/auth"
	"convention.ninja/internal/inventory/business"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	inventoryGrp := grp.Group("/orgs/:orgId/inventory")
	inventoryGrp.Use(auth.NewUserRequired())
	modelsGrp := inventoryGrp.Group("/models")
	setupModels(modelsGrp)

	mfgGrp := inventoryGrp.Group("/manufacturers")

	setupMfg(mfgGrp)

	categoriesGrp := inventoryGrp.Group("/categories")

	setupCategories(categoriesGrp)

	assetsGrp := inventoryGrp.Group("/assets")

	setupAssets(assetsGrp)

	manifestsGrp := inventoryGrp.Group("/manifests")

	setupManifests(manifestsGrp)
}

func setupManifests(manifestsGrp fiber.Router) {
	manifestsGrp.Get("/", business.GetManifests)

	manifestsGrp.Post("/", business.CreateManifest)

	manifestsGrp.Get("/:manifestId", business.GetManifest)

	manifestsGrp.Put("/:manifestId", business.UpdateManifest)

	manifestsGrp.Post("/:manifestId/ship", business.ShipManifest)

	manifestsGrp.Delete("/:manifestId/ship", business.UnshipManifest)

	manifestsGrp.Delete("/:manifestId", business.DeleteManifest)

	manifestsGrp.Get("/:manifestId/assets", business.GetManifestAssets)

	manifestsGrp.Post("/:manifestId/assets", business.AddAssetToManifest)

	manifestsGrp.Put("/:manifestId/assets", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Delete("/:manifestId/assets/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}

func setupAssets(assetsGrp fiber.Router) {
	assetsGrp.Get("/", business.GetAssets)

	assetsGrp.Post("/", business.CreateAsset)

	assetsGrp.Get("/:assetId", business.GetAsset)

	assetsGrp.Patch("/:assetId", business.UpdateAsset)

	assetsGrp.Get("/:assetId/barcodes", business.GetAssetBarcodes)

	assetsGrp.Post("/:assetId/barcodes", business.CreateAssetBarcode)

	assetsGrp.Delete("/:assetId/barcodes/:barcodeId", business.DeleteAssetBarcode)

	assetsGrp.Delete("/:assetId", business.DeleteAsset)

	assetsGrp.Get("/barcode/:barcode", business.GetAssetByBarcode)
}

func setupCategories(categoriesGrp fiber.Router) {
	categoriesGrp.Get("/", business.GetCategories)

	categoriesGrp.Post("/", business.CreateCategory)

	categoriesGrp.Get("/:categoryId", business.GetCategory)

	categoriesGrp.Put("/:categoryId", business.UpdateCategory)

	categoriesGrp.Delete("/:categoryId", business.DeleteCategory)
}

func setupMfg(mfgGrp fiber.Router) {
	mfgGrp.Get("/", business.GetManufacturers)

	mfgGrp.Post("/", business.CreateManufacturer)

	mfgGrp.Get("/:mfgId", business.GetManufacturer)

	mfgGrp.Put("/:mfgId", business.UpdateManufacturer)

	mfgGrp.Delete("/:mfgId", business.DeleteManufacturer)
}

func setupModels(modelsGrp fiber.Router) {
	modelsGrp.Get("/", business.GetModels)

	modelsGrp.Post("/", business.CreateModel)

	modelsGrp.Get("/:modelId", business.GetModel)

	modelsGrp.Put("/:modelId", business.UpdateModel)

	modelsGrp.Delete("/:modelId", business.DeleteModel)
}
