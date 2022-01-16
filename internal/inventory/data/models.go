package data

import (
	"convention.ninja/internal/data"
	data2 "convention.ninja/internal/organizations/data"
)

type Manufacturer struct {
	data.SnowflakeModel
	Name           string
	OrganizationId int64
	Organization   data2.Organization
}

type Category struct {
	data.SnowflakeModel
	Name           string
	OrganizationId int64
	Organization   data2.Organization
}

type Model struct {
	data.SnowflakeModel
	Name           string
	ManufacturerId int64
	Manufacturer   Manufacturer
	CategoryId     int64
	Category       Category
	OrganizationId int64
	Organization   data2.Organization
}

type AssetTag struct {
	data.SnowflakeModel
	TagId          string
	AssetId        int64
	Asset          Asset
	OrganizationId int64
}

type Asset struct {
	data.SnowflakeModel
	ModelId        int64
	Model          Model
	SerialNumber   string
	OrganizationId int64
	Organization   data2.Organization
	AssetTags      []AssetTag
}
