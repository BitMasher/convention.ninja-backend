package data

import (
	"convention.ninja/internal/data"
	data2 "convention.ninja/internal/organizations/data"
)

type Manufacturer struct {
	data.SnowflakeModel
	Name           string             `json:"name"`
	OrganizationId int64              `json:"organizationId,string"`
	Organization   data2.Organization `json:"-"`
}

type Category struct {
	data.SnowflakeModel
	Name           string             `json:"name"`
	OrganizationId int64              `json:"organizationId,string"`
	Organization   data2.Organization `json:"-"`
}

type Model struct {
	data.SnowflakeModel
	Name           string             `json:"name"`
	ManufacturerId int64              `json:"manufacturerId,string"`
	Manufacturer   Manufacturer       `json:"manufacturer,omitempty"`
	CategoryId     int64              `json:"categoryId,string"`
	Category       Category           `json:"category,omitempty"`
	OrganizationId int64              `json:"organizationId,string"`
	Organization   data2.Organization `json:"-"`
}

type AssetTag struct {
	data.SnowflakeModel
	TagId          string `json:"tagId"`
	AssetId        int64  `json:"assetId,string"`
	Asset          Asset  `json:"asset,omitempty"`
	OrganizationId int64  `json:"organizationId,string"`
}

type Asset struct {
	data.SnowflakeModel
	ModelId        int64              `json:"modelId,string"`
	Model          Model              `json:"model,omitempty"`
	SerialNumber   string             `json:"serialNumber"`
	OrganizationId int64              `json:"organizationId,string"`
	Organization   data2.Organization `json:"-"`
	AssetTags      []AssetTag         `json:"assetTags,omitempty"`
}

func Migrate() {
	_ = data.GetConn().AutoMigrate(&Manufacturer{})
	_ = data.GetConn().AutoMigrate(&Category{})
	_ = data.GetConn().AutoMigrate(&Model{})
	_ = data.GetConn().AutoMigrate(&Asset{})
	_ = data.GetConn().AutoMigrate(&AssetTag{})
}
