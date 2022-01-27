package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

func GetCategoriesByOrganization(orgId int64) (*[]Category, error) {
	rows, err := data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.categories where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	categories := make([]Category, 0)
	for rows.Next() {
		cat := Category{}
		if err = rows.Scan(&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, cat)
	}
	if err = rows.Err(); err != nil {
		return nil, err
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
	return false, errors.New("expected one row but got zero")
}

func CategoryExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.categories where id = ? and organization_id = ? and deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.categories where id = ? and deleted_at is null", id)
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
	defer rows.Close()
	if rows.Next() {
		var cat Category
		if err = rows.Scan(&cat.ID, &cat.Name, &cat.OrganizationId, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		return &cat, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
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
