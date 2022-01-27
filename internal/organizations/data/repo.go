package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"strings"
	"time"
)

func GetOrganizationsByOwner(ownerId int64) (*[]Organization, error) {
	// TODO implement get by association/permission instead of owner
	rows, err := data.GetConn().Query("select id, name, normalized_name, owner_id, created_at, updated_at, deleted_at from ods.organizations where owner_id = ? and deleted_at is null", ownerId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Organization
	results = make([]Organization, 0)
	for rows.Next() {
		org := Organization{}
		if err = rows.Scan(&org.ID, &org.Name, &org.NormalizedName, &org.OwnerId, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt); err != nil {
			return nil, err
		}
		results = append(results, org)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &results, nil
}

func OrganizationNameExists(name string) (bool, error) {
	rows, err := data.GetConn().Query("select count(id) from ods.organizations where normalized_name = ? and deleted_at is null", strings.ToLower(name))
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

func CreateOrganization(org *Organization) error {
	if org.OwnerId == 0 {
		return errors.New("owner id is a required property of an organization")
	}
	org.ID = snowflake.GetNode().Generate().Int64()
	org.NormalizedName = strings.ToLower(org.Name)
	org.CreatedAt = time.Now()
	org.UpdatedAt = org.CreatedAt
	org.DeletedAt = sql.NullTime{}
	res, err := data.GetConn().Exec("insert into ods.organizations (id, name, normalized_name, owner_id, created_at, updated_at) values(?, ?, ?, ?, ?, ?)", org.ID, org.Name, org.NormalizedName, org.OwnerId, org.CreatedAt, org.UpdatedAt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return errors.New("expected 1 affected row got different")
	}
	return nil
}

func GetOrganizationById(id int64) (*Organization, error) {
	rows, err := data.GetConn().Query("select id, name, normalized_name, owner_id, created_at, updated_at, deleted_at from ods.organizations where id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		var org Organization
		if err = rows.Scan(&org.ID, &org.Name, &org.NormalizedName, &org.OwnerId, &org.CreatedAt, &org.UpdatedAt, &org.DeletedAt); err != nil {
			return nil, err
		}
		return &org, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func UpdateOrganization(org *Organization) error {
	org.UpdatedAt = time.Now()
	if org.DeletedAt.Valid {
		return errors.New("cannot update organization that is already deleted")
	}
	org.NormalizedName = strings.ToLower(org.Name)
	res, err := data.GetConn().Exec("update ods.organizations set name = ?, normalized_name = ?, updated_at = ? where id = ? and deleted_at is null limit 1", org.Name, org.NormalizedName, org.UpdatedAt, org.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return errors.New("expected 1 affected rows got different")
	}
	return nil
}

func DeleteOrganization(org *Organization) error {
	if org.DeletedAt.Valid {
		return errors.New("cannot delete organization that is already deleted")
	}
	org.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.organizations set deleted_at = ? where id = ? and deleted_at is null limit 1", org.DeletedAt, org.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return errors.New("expected 1 row affected got different")
	}
	return nil
}
