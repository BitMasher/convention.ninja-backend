package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

func GetManufacturersByOrganization(orgId int64) (*[]Manufacturer, error) {
	rows, err := data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	manufacturers := make([]Manufacturer, 0)
	for rows.Next() {
		mfg := Manufacturer{}
		if err = rows.Scan(&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
			return nil, err
		}
		manufacturers = append(manufacturers, mfg)
	}
	if err = rows.Err(); err != nil {
		return nil, err
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

func ManufacturerExistsById(id int64, orgId ...int64) (bool, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select count(id) from ods.manufacturers where id = ? and organization_id = ? and deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select count(id) from ods.manufacturers where id = ? and deleted_at is null", id)
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

func GetManufacturerById(id int64, orgId ...int64) (*Manufacturer, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where id = ? and organization_id = ? and deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select id, name, organization_id, created_at, updated_at from ods.manufacturers where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var mfg Manufacturer
		if err = rows.Scan(&mfg.ID, &mfg.Name, &mfg.OrganizationId, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
			return nil, err
		}
		return &mfg, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
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
