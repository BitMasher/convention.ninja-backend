package inventory

import (
	"convention.ninja/internal/auth"
	"errors"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(grp fiber.Router, _ *auth.IgnoreListBuilder) {
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
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Get("/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Patch("/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Put("/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Delete("/:assetId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	assetsGrp.Get("/barcode/:barcode", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}

func setupCategories(categoriesGrp fiber.Router) {
	categoriesGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	categoriesGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	categoriesGrp.Get("/:categoryId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	categoriesGrp.Put("/:categoryId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	categoriesGrp.Delete("/:categoryId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}

func setupMfg(mfgGrp fiber.Router) {
	mfgGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	mfgGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	mfgGrp.Get("/:mfgId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	mfgGrp.Put("/:mfgId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	mfgGrp.Delete("/:mfgId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}

func setupModels(modelsGrp fiber.Router) {
	modelsGrp.Get("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	modelsGrp.Post("/", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	modelsGrp.Get("/:modelId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	modelsGrp.Put("/:modelId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})

	modelsGrp.Delete("/:modelId", func(c *fiber.Ctx) error {
		return errors.New("not implemented") // TODO
	})
}
