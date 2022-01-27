package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func GetAssetsByOrganization(orgId int64) (*[]Asset, error) {
	rows, err := data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assets := make([]Asset, 0)
	for rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &assets, nil
}

func GetAssetsExpandedByOrganization(orgId int64) (*[]Asset, error) {
	rows, err := data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at, m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m.deleted_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.assets a inner join ods.models m on m.id = a.model_id inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where a.organization_id = ? and a.deleted_at is null order by a.id", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assets := make([]Asset, 0)
	ids := make([]interface{}, 0)
	for rows.Next() {
		asset := Asset{}
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt,
			&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt, &model.DeletedAt,
			&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt, &mfg.DeletedAt,
			&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		asset.Model = model
		asset.AssetTags = make([]AssetTag, 0)
		assets = append(assets, asset)
		ids = append(ids, asset.ID)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	if len(ids) > 0 {
		rows2, err := data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id in (?"+strings.Repeat(",?", len(ids)-1)+") and deleted_at is null order by asset_id", ids...)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		assetIdx := 0
		maxIdx := len(assets)
		for rows2.Next() {
			tag := AssetTag{}
			if err = rows2.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
				return nil, err
			}
			if assets[assetIdx].ID != tag.AssetId {
				assetIdx++
			}
			if assetIdx < maxIdx && assets[assetIdx].ID != tag.AssetId {
				for i := assetIdx; i < len(assets); i++ {
					if assets[i].ID == tag.AssetId {
						assetIdx = i
					}
				}
			}
			if assetIdx < maxIdx && assets[assetIdx].ID == tag.AssetId {
				assets[assetIdx].AssetTags = append(assets[assetIdx].AssetTags, tag)
			}
		}
		if err = rows2.Err(); err != nil {
			return nil, err
		}
	}
	return &assets, nil
}

func CreateAsset(asset *Asset) error {
	asset.ID = snowflake.GetNode().Generate().Int64()
	asset.CreatedAt = time.Now()
	asset.UpdatedAt = time.Now()
	asset.DeletedAt = sql.NullTime{}

	res, err := data.GetConn().Exec("insert into ods.assets (id, model_id, serial_number, organization_id, created_at, updated_at) values(?,?,?,?,?,?)", asset.ID, asset.ModelId, asset.SerialNumber, asset.OrganizationId, asset.CreatedAt, asset.UpdatedAt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one affected row got different")
	}
	return nil
}

func BulkCreateAssetTag(assetTags []AssetTag) error {
	values := make([]interface{}, 0)
	insert := make([]string, 0)
	for i := range assetTags {
		assetTags[i].ID = snowflake.GetNode().Generate().Int64()
		assetTags[i].CreatedAt = time.Now()
		assetTags[i].UpdatedAt = time.Now()
		assetTags[i].DeletedAt = sql.NullTime{}
		values = append(values, assetTags[i].ID, assetTags[i].TagId, assetTags[i].AssetId, assetTags[i].OrganizationId, assetTags[i].CreatedAt, assetTags[i].UpdatedAt)
		insert = append(insert, "(?,?,?,?,?,?)")
	}
	res, err := data.GetConn().Exec("insert into ods.asset_tags (id, tag_id, asset_id, organization_id, created_at, updated_at) values"+strings.Join(insert, ",")+"", values...)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != int64(len(assetTags)) {
		return errors.New("expected affected rows to match count of tags but got different")
	}
	return nil
}

func AssetTagsExistInOrg(orgId int64, assetTags []string) (bool, error) {
	if len(assetTags) == 0 {
		return false, nil
	}
	tagIds := make([]interface{}, 0, len(assetTags)+1)
	tagIds = append(tagIds, orgId)
	for _, tag := range assetTags {
		tagIds = append(tagIds, tag)
	}

	rows, err := data.GetConn().Query("select count(id) from ods.asset_tags where organization_id = ? and tag_id in (?"+strings.Repeat(",?", len(assetTags)-1)+") and deleted_at is null", tagIds...)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}
	return false, errors.New("expected one result got zero")
}

func AssetTagExistInOrg(orgId int64, assetTag string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.asset_tags where organization_id = ? and tag_id = ? and deleted_at is null", orgId, assetTag)
	if err != nil {
		return false, nil
	}
	defer rows.Close()
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}
	return false, errors.New("expected one result got zero")
}

func GetAssetById(id int64, orgId ...int64) (*Asset, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func GetAssetExpandedById(id int64, orgId ...int64) (*Asset, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at, m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m.deleted_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.assets a inner join ods.models m on m.id = a.model_id inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where a.organization_id = ? and a.id = ? and a.deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at, m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m.deleted_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.assets a inner join ods.models m on m.id = a.model_id inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where a.id = ? and a.deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		asset := Asset{}
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt,
			&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt, &model.DeletedAt,
			&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt, &mfg.DeletedAt,
			&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		asset.Model = model
		asset.AssetTags = make([]AssetTag, 0)

		rows2, err := data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id = ? and deleted_at is null", asset.ID)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			tag := AssetTag{}
			if err = rows2.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
				return nil, err
			}
			asset.AssetTags = append(asset.AssetTags, tag)
		}
		if err = rows2.Err(); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func UpdateAsset(asset *Asset) error {
	if asset.DeletedAt.Valid {
		return errors.New("cannot update asset that has already been deleted")
	}
	asset.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.assets set serial_number = ?, model_id = ?, updated_at = ? where id = ? and organization_id = ? limit 1", asset.SerialNumber, asset.ModelId, asset.UpdatedAt, asset.ID, asset.OrganizationId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one affected row got different")
	}
	return nil
}

func GetAssetTagsByAssetId(assetId int64, orgId ...int64) (*[]AssetTag, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where organization_id = ? and asset_id = ? and deleted_at is null", orgId[0], assetId)
	} else {
		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id = ? and deleted_at is null", assetId)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	assetTags := make([]AssetTag, 0)
	for rows.Next() {
		assetTag := AssetTag{}
		if err = rows.Scan(&assetTag.ID, &assetTag.TagId, &assetTag.AssetId, &assetTag.OrganizationId, &assetTag.CreatedAt, &assetTag.UpdatedAt); err != nil {
			return nil, err
		}
		assetTags = append(assetTags, assetTag)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &assetTags, nil
}

func AssetExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.assets where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.assets where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return false, err
	}
	defer rows.Close()
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}
	return false, errors.New("expected one result got zero")
}

func CreateAssetTag(tag *AssetTag) error {
	tag.ID = snowflake.GetNode().Generate().Int64()
	tag.CreatedAt = time.Now()
	tag.UpdatedAt = time.Now()
	tag.DeletedAt = sql.NullTime{}
	res, err := data.GetConn().Exec("insert into ods.asset_tags (id, tag_id, asset_id, organization_id, created_at, updated_at) values(?, ?, ?, ?, ?, ?)", tag.ID, tag.TagId, tag.AssetId, tag.OrganizationId, tag.CreatedAt, tag.UpdatedAt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one affected row got different")
	}
	return nil
}

func DeleteAssetTag(tag *AssetTag) error {
	if tag.DeletedAt.Valid {
		return errors.New("cannot delete asset tag that has already been deleted")
	}
	tag.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.asset_tags set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", tag.DeletedAt, tag.ID, tag.OrganizationId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected got different")
	}
	return nil
}

func GetAssetTagById(id int64, orgId ...int64) (*AssetTag, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		tag := AssetTag{}
		if err = rows.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, err
		}
		return &tag, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func DeleteAsset(asset *Asset) error {
	if asset.DeletedAt.Valid {
		return errors.New("cannot delete asset that has already been deleted")
	}
	asset.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.assets set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", asset.DeletedAt, asset.ID, asset.OrganizationId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected got different")
	}
	return nil
}

func DeleteAssetTagsByAssetId(assetId int64, orgId ...int64) error {
	var err error
	dt := sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	if len(orgId) > 0 {
		_, err = data.GetConn().Exec("update ods.asset_tags set deleted_at = ? where asset_id = ? and organization_id = ? and deleted_at is null", dt, assetId, orgId[0])
	} else {
		_, err = data.GetConn().Exec("update ods.asset_tags set deleted_at = ? where asset_id = ? and deleted_at is null", dt, assetId)
	}
	if err != nil {
		return err
	}
	return nil
}

func GetAssetByTag(tag string, orgId int64) (*Asset, error) {
	rows, err := data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at from ods.assets a inner join ods.asset_tags t on t.asset_id = a.id and t.deleted_at is null where a.organization_id = ? and t.tag_id = ? and a.deleted_at is null", orgId, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func GetAssetExpandedByTag(tag string, orgId int64) (*Asset, error) {
	rows, err := data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at, m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m.deleted_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.assets a inner join ods.models m on m.id = a.model_id inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id inner join ods.asset_tags t on t.asset_id = a.id and t.deleted_at is null where a.organization_id = ? and t.tag_id = ? and a.deleted_at is null", orgId, tag)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		asset := Asset{}
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt,
			&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt, &model.DeletedAt,
			&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt, &mfg.DeletedAt,
			&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		asset.Model = model
		asset.AssetTags = make([]AssetTag, 0)

		rows2, err := data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id = ? and deleted_at is null", asset.ID)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			tag := AssetTag{}
			if err = rows2.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
				return nil, err
			}
			asset.AssetTags = append(asset.AssetTags, tag)
		}
		if err = rows2.Err(); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}
