package data

import (
	"convention.ninja/internal/data"
	data2 "convention.ninja/internal/organizations/data"
	"database/sql"
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
	Asset          Asset  `json:"-"`
	OrganizationId int64  `json:"organizationId,string"`
}

type Asset struct {
	data.SnowflakeModel
	ModelId        int64              `json:"modelId,string"`
	Model          Model              `json:"model,omitempty"`
	SerialNumber   string             `json:"serialNumber"`
	RoomId         string             `json:"roomId"`
	OrganizationId int64              `json:"organizationId,string"`
	Organization   data2.Organization `json:"-"`
	AssetTags      []AssetTag         `json:"assetTags,omitempty"`
}

type Manifest struct {
	data.SnowflakeModel
	RoomId                   string          `json:"locationId"`
	ResponsibleUserId        sql.NullInt64   `json:"responsibleUserId"`
	ResponsibleExternalParty sql.NullString  `json:"responsibleExternalParty"`
	ShipDate                 sql.NullTime    `json:"shipDate"`
	OrganizationId           int64           `json:"organizationId"`
	CreatorId                int64           `json:"creatorId"`
	Entries                  []ManifestEntry `json:"entries"`
}

type ManifestEntry struct {
	data.SnowflakeModel
	ManifestId     int64 `json:"manifestId"`
	AssetId        int64 `json:"assetId"`
	OrganizationId int64 `json:"organizationId"`
}

type ExternalPartyIdentifier struct {
	data.SnowflakeModel
	Fields         string `json:"fields"`
	OrganizationId int64  `json:"organizationId,string"`
}
