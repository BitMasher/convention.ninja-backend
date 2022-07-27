package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

func GetModelsByOrganization(orgId int64) (*[]Model, error) {
	rows, err := data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	models := make([]Model, 0)
	for rows.Next() {
		model := Model{}
		if err = rows.Scan(&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		models = append(models, model)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &models, nil
}

func GetModelsExpandedByOrganization(orgId int64) (*[]Model, error) {
	rows, err := data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.organization_id = ? and m.deleted_at is null order by m.id", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &models, nil
}

func ModelExistsInOrg(orgId int64, name string, mfgId int64) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.models where organization_id = ? and name = ? and manufacturer_id = ? and deleted_at is null", orgId, name, mfgId)
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
		rows, err = data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where id = ? and organization_id = ? and deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select id, name, manufacturer_id, category_id, organization_id, created_at, updated_at from ods.models where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		model := Model{}
		if err = rows.Scan(&model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.OrganizationId, &model.CreatedAt, &model.UpdatedAt); err != nil {
			return nil, err
		}
		return &model, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func GetModelExpandedById(id int64, orgId ...int64) (*Model, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.id = ? and m.organization_id = ? and m.deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select m.id, m.name, m.manufacturer_id, m.category_id, m.organization_id, m.created_at, m.updated_at, m2.id, m2.name, m2.organization_id, m2.created_at, m2.updated_at, m2.deleted_at, c.id, c.name, c.organization_id, c.created_at, c.updated_at, c.deleted_at from ods.models m inner join ods.manufacturers m2 on m2.id = m.manufacturer_id inner join ods.categories c on c.id = m.category_id where m.id = ? and m.deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
	if err = rows.Err(); err != nil {
		return nil, err
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
