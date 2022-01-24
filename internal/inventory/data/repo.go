package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func GetCategoriesByOrganization(orgId int64) (*[]Category, error) {
	rows, err := data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.categories where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	categories := make([]Category, 0)
	for rows.Next() {
		cat := Category{}
		if err = rows.Scan(&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	return &categories, nil
}

func CreateCategory(cat *Category) error {
	cat.ID = snowflake.GetNode().Generate().Int64()
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()
	cat.DeletedAt = sql.NullTime{}

	res, err := data.GetConn().Exec("insert into ods.categories (id, name, organization_id, created_at, updated_at) values (?,?,?,?,?)", cat.ID, cat.Name, cat.OrganizationId, cat.CreatedAt, cat.UpdatedAt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return errors.New("expected 1 affected row got different")
	}
	return nil
}

func CategoryExistsInOrg(orgId int64, name string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.categories where organization_id = ? and name = ? and deleted_at is null", orgId, name)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one row but got zero")
}

func CategoryExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.categories where id = ? and organization_id = ? and deleted_at is null", id, orgId)
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.categories where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func GetCategoryById(id int64, orgId ...int64) (*Category, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.categories where id = ? and organization_id = ? and deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.categories where id = ? and deleted_at is null")
	}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var cat Category
		if err = rows.Scan(&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		return &cat, nil
	}
	return nil, nil
}

func UpdateCategory(cat *Category) error {
	if cat.DeletedAt.Valid == true {
		return errors.New("cannot update category that has already been deleted")
	}
	cat.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.categories set name = ?, updated_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", cat.Name, cat.UpdatedAt, cat.ID, cat.OrganizationId)
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

func DeleteCategory(cat *Category) error {
	if cat.DeletedAt.Valid == true {
		return errors.New("cannot delete category that has already been deleted")
	}
	cat.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.categories set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", cat.DeletedAt, cat.ID, cat.OrganizationId)
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

func GetManufacturersByOrganization(orgId int64) (*[]Manufacturer, error) {
	rows, err := data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	manufacturers := make([]Manufacturer, 0)
	for rows.Next() {
		mfg := Manufacturer{}
		if err = rows.Scan(&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, mfg)
	}
	return &manufacturers, nil
}

func CreateManufacturer(mfg *Manufacturer) error {
	mfg.ID = snowflake.GetNode().Generate().Int64()
	mfg.CreatedAt = time.Now()
	mfg.UpdatedAt = time.Now()
	mfg.DeletedAt = sql.NullTime{}
	res, err := data.GetConn().Exec("insert into ods.manufacturers (id, name, organization_id, created_at, updated_at) values(?,?,?,?,?)", mfg.ID, mfg.Name, mfg.OrganizationId, mfg.CreatedAt, mfg.UpdatedAt)
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

func ManufacturerExistsInOrg(orgId int64, name string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.manufacturers where organization_id = ? and name = ? and deleted_at is null", orgId, name)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func ManufacturerExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.manufacturers where id = ? and organization_id = ? and deleted_at is null", id, orgId)
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.manufacturers where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func GetManufacturerById(id int64, orgId ...int64) (*Manufacturer, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where id = ? and organization_id = ? and deleted_at is null", id, orgId)
	} else {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		var mfg Manufacturer
		if err = rows.Scan(&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
			return nil, err
		}
		return &mfg, nil
	}
	return nil, nil
}

func UpdateManufacturer(mfg *Manufacturer) error {
	if mfg.DeletedAt.Valid {
		return errors.New("cannot update manufacturer that has already been deleted")
	}
	mfg.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.manufacturers set name = ?, updated_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", mfg.Name, mfg.UpdatedAt, mfg.ID, mfg.OrganizationId)
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

func DeleteManufacturer(mfg *Manufacturer) error {
	if mfg.DeletedAt.Valid {
		return errors.New("cannot delete manufacturer that has already been deleted")
	}
	mfg.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.manufacturers set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", mfg.DeletedAt, mfg.ID, mfg.OrganizationId)
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

func GetModelsByOrganization(orgId int64) (*[]Model, error) {
	rows, err := data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	models := make([]Model, 0)
	for rows.Next() {
		model := Model{}
		if err = rows.Scan(&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	return &models, nil
}

func GetModelsExpandedByOrganization(orgId int64) (*[]Model, error) {
	rows, err := data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.organization_id = ? and m.deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	models := make([]Model, 0)
	for rows.Next() {
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt,
			&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt, &mfg.DeletedAt,
			&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		models = append(models, model)
	}
	return &models, nil
}

func ModelExistsInOrg(orgId int64, name string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.models where organization_id = ? and name = ? and deleted_at is null", orgId, name)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func ModelExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.models where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.models where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func CreateModel(model *Model) error {
	model.ID = snowflake.GetNode().Generate().Int64()
	model.CreatedAt = time.Now()
	model.UpdatedAt = time.Now()
	model.DeletedAt = sql.NullTime{}

	res, err := data.GetConn().Exec("insert into ods.models (id, name, manufacturer_id, category_id, organization_id, created_at, updated_at) values(?,?,?,?,?,?,?)", model.ID, model.Name, model.ManufacturerId, model.CategoryId, model.OrganizationId, model.CreatedAt, model.UpdatedAt)
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

func GetModelById(id int64, orgId ...int64) (*Model, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where id = ? and organization_id = ? and deleted_at is null", id, orgId)
	} else {
		rows, err = data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		model := Model{}
		if err = rows.Scan(&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		return &model, nil
	}
	return nil, nil
}

func GetModelExpandedById(id int64, orgId ...int64) (*Model, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.id = ? and m.organization_id = ? and m.deleted_at is null", id, orgId)
	} else {
		rows, err = data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.id = ? and m.deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt,
			&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt, &mfg.DeletedAt,
			&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt, &cat.DeletedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		return &model, nil
	}
	return nil, nil
}

func UpdateModel(model *Model) error {
	if model.DeletedAt.Valid {
		return errors.New("cannot update model that has already been deleted")
	}
	model.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.models set name = ?, manufacturer_id = ?, category_id = ?, updated_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", model.Name, model.ManufacturerId, model.CategoryId, model.UpdatedAt, model.ID, model.OrganizationId)
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

func DeleteModel(model *Model) error {
	if model.DeletedAt.Valid {
		return errors.New("cannot delete model that has already been deleted")
	}
	model.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.models set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null", model.DeletedAt, model.ID, model.OrganizationId)
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

func GetAssetsByOrganization(orgId int64) (*[]Asset, error) {
	rows, err := data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	assets := make([]Asset, 0)
	for rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}
	return &assets, nil
}

func GetAssetsExpandedByOrganization(orgId int64) (*[]Asset, error) {
	rows, err := data.GetConn().Query("select a.id, a.model_id, a.serial_number, a.organization_id, a.created_at, a.updated_at, m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m.deleted_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.assets a inner join ods.models m on m.id = a.model_id inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where a.organization_id = ? and a.deleted_at is null order by a.id", orgId)
	if err != nil {
		return nil, err
	}
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
	if len(ids) > 0 {
		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id in (?"+strings.Repeat(",?", len(ids)-1)+") and deleted_at is null order by asset_id", ids...)
		if err != nil {
			return nil, err
		}
		assetIdx := 0
		for rows.Next() {
			tag := AssetTag{}
			if err = rows.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
				return nil, err
			}
			if assets[assetIdx].ID != tag.ID {
				assetIdx++
			}
			if assets[assetIdx].ID != tag.ID {
				for i := assetIdx; i < len(assets); i++ {
					if assets[i].ID == tag.ID {
						assetIdx = i
					}
				}
			}
			if assets[assetIdx].ID == tag.ID {
				assets[assetIdx].AssetTags = append(assets[assetIdx].AssetTags, tag)
			}
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
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func AssetTagExistInOrg(orgId int64, assetTag string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.asset_tags where organization_id = ? and tag_id = ? and deleted_at is null", orgId, assetTag)
	if err != nil {
		return false, nil
	}
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
	}
	return false, errors.New("expected one result got zero")
}

func GetAssetById(id int64, orgId ...int64) (*Asset, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where organization_id = ? and id = ? and deleted_at is null", orgId, id)
	} else {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	if rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	return nil, nil
}

func GetAssetExpandedById(id int64, orgId ...int64) (*Asset, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where organization_id = ? and id = ? and deleted_at is null", orgId, id)
	} else {
		rows, err = data.GetConn().Query("select id, model_id, serial_number, organization_id, created_at, updated_at from ods.assets where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
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

		rows, err = data.GetConn().Query("select id, tag_id, asset_id, organization_id, created_at, updated_at from ods.asset_tags where asset_id = ? and deleted_at is null", asset.ID)
		if err != nil {
			return nil, err
		}
		for rows.Next() {
			tag := AssetTag{}
			if err = rows.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
				return nil, err
			}
			asset.AssetTags = append(asset.AssetTags, tag)
		}
		return &asset, nil
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
	assetTags := make([]AssetTag, 0)
	for rows.Next() {
		assetTag := AssetTag{}
		if err = rows.Scan(&assetTag.ID, &assetTag.TagId, &assetTag.AssetId, &assetTag.OrganizationId, &assetTag.CreatedAt, &assetTag.UpdatedAt); err != nil {
			return nil, err
		}
		assetTags = append(assetTags, assetTag)
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
	if rows.Next() {
		var count int64
		if err = rows.Scan(&count); err != nil {
			return false, err
		}
		return count > 0, nil
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
	res, err := data.GetConn().Exec("update ods.asset_tags set deleted_at = ? where id = ? and organization_id = ? limit 1", tag.DeletedAt, tag.ID, tag.OrganizationId)
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
	if rows.Next() {
		tag := AssetTag{}
		if err = rows.Scan(&tag.ID, &tag.TagId, &tag.AssetId, &tag.OrganizationId, &tag.CreatedAt, &tag.UpdatedAt); err != nil {
			return nil, err
		}
		return &tag, nil
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
	res, err := data.GetConn().Exec("update ods.assets set deleted_at = ? where id = ? and organization_id = ? and deleted_at is not null limit 1", asset.DeletedAt, asset.ID, asset.OrganizationId)
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
	if rows.Next() {
		asset := Asset{}
		if err = rows.Scan(&asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.OrganizationId, &asset.CreatedAt, &asset.UpdatedAt); err != nil {
			return nil, err
		}
		return &asset, nil
	}
	return nil, nil
}
