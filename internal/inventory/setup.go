package inventory

import (
	"convention.ninja/internal/inventory/business"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router) {
	inventoryGrp := grp.Group("/orgs/:orgId/inventory")
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
	manifestsGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Get("/:manifestId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Put("/:manifestId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Patch("/:manifestId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Delete("/:manifestId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Get("/:manifestId/assets", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Post("/:manifestId/assets", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Put("/:manifestId/assets", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	manifestsGrp.Delete("/:manifestId/assets/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}

func setupAssets(assetsGrp fiber.Router) {
	assetsGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetAssets(c)
	})

	assetsGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateAsset(c)
	})

	assetsGrp.Get("/:assetId", func(c *fiber.Ctx) error {
		return business.GetAsset(c)
	})

	assetsGrp.Patch("/:assetId", func(c *fiber.Ctx) error {
		return business.UpdateAsset(c)
	})

	assetsGrp.Get("/:assetId/barcodes", func(c *fiber.Ctx) error {
		return business.GetAssetBarcodes(c)
	})

	assetsGrp.Post("/:assetId/barcodes", func(c *fiber.Ctx) error {
		return business.CreateAssetBarcode(c)
	})

	assetsGrp.Delete("/:assetId/barcodes/:barcodeId", func(c *fiber.Ctx) error {
		return business.DeleteAssetBarcode(c)
	})

	assetsGrp.Delete("/:assetId", func(c *fiber.Ctx) error {
		return business.DeleteAsset(c)
	})

	assetsGrp.Get("/barcode/:barcode", func(c *fiber.Ctx) error {
		return business.GetAssetByBarcode(c)
	})
}

func setupCategories(categoriesGrp fiber.Router) {
	categoriesGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetCategories(c)
	})

	categoriesGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateCategory(c)
	})

	categoriesGrp.Get("/:categoryId", func(c *fiber.Ctx) error {
		return business.GetCategory(c)
	})

	categoriesGrp.Put("/:categoryId", func(c *fiber.Ctx) error {
		return business.UpdateCategory(c)
	})

	categoriesGrp.Delete("/:categoryId", func(c *fiber.Ctx) error {
		return business.DeleteCategory(c)
	})
}

func setupMfg(mfgGrp fiber.Router) {
	mfgGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetManufacturers(c)
	})

	mfgGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateManufacturer(c)
	})

	mfgGrp.Get("/:mfgId", func(c *fiber.Ctx) error {
		return business.GetManufacturer(c)
	})

	mfgGrp.Put("/:mfgId", func(c *fiber.Ctx) error {
		return business.UpdateManufacturer(c)
	})

	mfgGrp.Delete("/:mfgId", func(c *fiber.Ctx) error {
		return business.DeleteManufacturer(c)
	})
}

func setupModels(modelsGrp fiber.Router) {
	modelsGrp.Get("/", func(c *fiber.Ctx) error {
		return business.GetModels(c)
	})

	modelsGrp.Post("/", func(c *fiber.Ctx) error {
		return business.CreateModel(c)
	})

	modelsGrp.Get("/:modelId", func(c *fiber.Ctx) error {
		return business.GetModel(c)
	})

	modelsGrp.Put("/:modelId", func(c *fiber.Ctx) error {
		return business.UpdateModel(c)
	})

	modelsGrp.Delete("/:modelId", func(c *fiber.Ctx) error {
		return business.DeleteModeL(c)
	})
}
