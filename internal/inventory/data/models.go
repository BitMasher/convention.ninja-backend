package data

import (
	"convention.ninja/internal/data"
	data2 "convention.ninja/internal/organizations/data"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
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
	RoomId                   sql.NullString  `json:"locationId"`
	ResponsibleUserId        sql.NullInt64   `json:"responsibleUserId"`
	ResponsibleExternalParty ExternalParty   `json:"responsibleExternalParty"`
	ShipDate                 sql.NullTime    `json:"shipDate"`
	OrganizationId           int64           `json:"organizationId,string"`
	CreatorId                int64           `json:"creatorId,string"`
	Entries                  []ManifestEntry `json:"entries"`
}

type ManifestEntry struct {
	data.SnowflakeModel
	ManifestId     int64 `json:"manifestId,string"`
	AssetId        int64 `json:"assetId,string"`
	Asset          Asset `json:"asset,omitempty"`
	OrganizationId int64 `json:"organizationId,string"`
}

type ExternalPartyIdentifier struct {
	data.SnowflakeModel
	Fields         string `json:"fields"`
	OrganizationId int64  `json:"organizationId,string"`
}

type ExternalParty struct {
	Name  string `json:"name"`
	Extra string `json:"extra"`
	Valid bool   `json:"valid"`
}

func (e *ExternalParty) Scan(src interface{}) error {
	var source []byte
	_m := ExternalParty{}

	switch src.(type) {
	case []uint8:
		source = src.([]uint8)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for ExternalParty")
	}
	err := json.Unmarshal(source, &_m)
	if err != nil {
		return err
	}
	_m.Valid = true
	*e = _m
	return nil
}

func (e ExternalParty) Value() (driver.Value, error) {
	if len(e.Name) == 0 && len(e.Extra) == 0 {
		return driver.Value(nil), nil
	}
	return json.Marshal(e)
}
