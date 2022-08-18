package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

func GetOpenManifestsByOrganization(orgId int64) (*[]Manifest, error) {
	rows, err := data.GetConn().Query("select id, room_id, responsible_external_party, ship_date, creator_id, organization_id, created_at, updated_at from ods.manifests m where m.organization_id = ? and deleted_at is null and ship_date is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	manifests := make([]Manifest, 0)
	for rows.Next() {
		manifest := Manifest{}
		if err = rows.Scan(&manifest.ID, &manifest.RoomId, &manifest.ResponsibleExternalParty, &manifest.ShipDate, &manifest.CreatorId, &manifest.OrganizationId, &manifest.CreatedAt, &manifest.UpdatedAt); err != nil {
			return nil, err
		}
		manifests = append(manifests, manifest)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &manifests, nil
}

func CreateManifest(manifest *Manifest) error {
	manifest.ID = snowflake.GetNode().Generate().Int64()
	manifest.CreatedAt = time.Now()
	manifest.UpdatedAt = time.Now()
	manifest.DeletedAt = sql.NullTime{}

	res, err := data.GetConn().Exec("insert into ods.manifests (id, room_id, creator_id, organization_id, created_at, updated_at) values(?,?,?,?,?,?)", manifest.ID, manifest.RoomId, manifest.CreatorId, manifest.OrganizationId, manifest.CreatedAt, manifest.UpdatedAt)
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

func GetManifestById(id int64, orgId ...int64) (*Manifest, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select m.id, m.room_id, m.responsible_external_party, m.ship_date, m.creator_id, m.organization_id, m.created_at, m.updated_at from ods.manifests m where m.id = ? and m.organization_id = ? and m.deleted_at is null", id, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select m.id, m.room_id, m.responsible_external_party, m.ship_date, m.creator_id, m.organization_id, m.created_at, m.updated_at from ods.manifests m where m.id = ? and m.deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	manifest := Manifest{}
	if rows.Next() {
		if err = rows.Scan(&manifest.ID, &manifest.RoomId, &manifest.ResponsibleExternalParty, &manifest.ShipDate, &manifest.CreatorId, &manifest.OrganizationId, &manifest.CreatedAt, &manifest.UpdatedAt); err != nil {
			return nil, err
		}
		manifest.Entries = make([]ManifestEntry, 0)
		rows2, err := data.GetConn().Query("select id, manifest_id, asset_id, organization_id, created_at, updated_at from ods.manifest_entries where manifest_id = ? and organization_id = ? and deleted_at is null", manifest.ID, manifest.OrganizationId)
		if err != nil {
			return nil, err
		}
		defer rows2.Close()
		for rows2.Next() {
			manifestEntry := ManifestEntry{}
			if err = rows2.Scan(&manifestEntry.ID, &manifestEntry.ManifestId, &manifestEntry.AssetId, &manifestEntry.OrganizationId, &manifestEntry.CreatedAt, &manifestEntry.UpdatedAt); err != nil {
				return nil, err
			}
			manifest.Entries = append(manifest.Entries, manifestEntry)
		}
		if err = rows2.Err(); err != nil {
			return nil, err
		}
		return &manifest, nil
	}
	return nil, nil
}

func GetManifestExpandedById(id int64, orgId ...int64) (*Manifest, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, room_id, ship_date, creator_id, organization_id, created_at, updated_at from ods.manifests m where m.organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select id, room_id, ship_date, creator_id, organization_id, created_at, updated_at from ods.manifests m where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	manifest := Manifest{}
	if rows.Next() {
		if err = rows.Scan(&manifest.ID, &manifest.RoomId, &manifest.ShipDate, &manifest.CreatorId, &manifest.OrganizationId, &manifest.CreatedAt, &manifest.UpdatedAt); err != nil {
			return nil, err
		}
		manifest.Entries = make([]ManifestEntry, 0)
		rows2, err := data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.created_at, e.updated_at, a.id, a.model_id, a.serial_number, a.room_id, a.created_at, a.updated_at, m2.id, m2.name, m2.manufacturer_id, m2.category_id, m2.created_at, m2.updated_at, c.id, c.name, c.created_at, c.updated_at, m3.id, m3.name, m3.created_at, m3.updated_at from ods.manifest_entries e inner join ods.assets a on a.id = e.asset_id inner join ods.models m2 on m2.id = a.model_id inner join ods.categories c on c.id = m2.category_id inner join ods.manufacturers m3 on m3.id = m2.manufacturer_id where e.manifest_id = ? and e.organization_id = ? and e.deleted_at is null", manifest.ID, manifest.OrganizationId)
		if err != nil {
			return nil, err
		}
		defer rows2.Next()
		for rows2.Next() {
			entry := ManifestEntry{}
			asset := Asset{}
			model := Model{}
			mfg := Manufacturer{}
			cat := Category{}
			if err = rows.Scan(
				&entry.ID, &entry.ManifestId, &entry.AssetId, &entry.CreatedAt, &entry.UpdatedAt, &asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.RoomId, &asset.CreatedAt, &asset.UpdatedAt, &model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.CreatedAt, &model.UpdatedAt, &cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt, &mfg.ID, &mfg.Name, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
				return nil, err
			}
			model.Manufacturer = mfg
			model.Category = cat
			asset.Model = model
			assetTags, err := GetAssetTagsByAssetId(asset.ID, manifest.OrganizationId)
			if err != nil {
				return nil, err
			}
			asset.AssetTags = *assetTags
			entry.Asset = asset
			manifest.Entries = append(manifest.Entries, entry)
		}
		if err = rows2.Err(); err != nil {
			return nil, err
		}
		return &manifest, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func UpdateManifest(manifest *Manifest) error {
	if manifest.DeletedAt.Valid {
		return errors.New("cannot update manifest that has already been deleted")
	}
	manifest.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.manifests set room_id = ?, responsible_external_party = ?, updated_at = ? where id = ? and deleted_at is null", manifest.RoomId, manifest.ResponsibleExternalParty, manifest.UpdatedAt, manifest.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected but got different")
	}
	return nil
}

func ShipManifest(manifest *Manifest) error {
	if manifest.DeletedAt.Valid {
		return errors.New("cannot ship manifest that has already been deleted")
	}
	if manifest.ShipDate.Valid {
		return errors.New("cannot ship manifest that has already been shipped")
	}
	if !manifest.RoomId.Valid && !manifest.ResponsibleExternalParty.Valid {
		return errors.New("cannot ship manifest with no destination")
	}
	manifest.ShipDate = sql.NullTime{
		Valid: true,
		Time:  time.Now(),
	}
	manifest.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.assets a inner join ods.manifest_entries e on e.asset_id = a.id inner join ods.manifests m on m.id = e.manifest_id set a.room_id = coalesce(m.room_id,concat('manifest:', m.id)), a.updated_at = ? where a.deleted_at is null and e.deleted_at is null and e.manifest_id = ? and m.deleted_at is null", manifest.UpdatedAt, manifest.ID)
	if err != nil {
		return err
	}
	res, err = data.GetConn().Exec("update ods.manifests set ship_date = ?, updated_at = ? where id = ? and deleted_at is null", manifest.ShipDate, manifest.UpdatedAt, manifest.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected but got different")
	}
	return nil
}

func UnshipManifest(manifest *Manifest) error {
	// TODO implement
	return errors.New("not implemented")
}

func DeleteManifest(manifest *Manifest) error {
	if manifest.DeletedAt.Valid {
		return errors.New("cannot delete manifest that has already been deleted")
	}
	manifest.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.manifests m set m.deleted_at = ? where m.id = ? and m.deleted_at is null", manifest.DeletedAt, manifest.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected but got different")
	}
	_, err = data.GetConn().Exec("update ods.manifest_entries e set e.deleted_at = ? where e.manifest_id = ? and e.deleted_at is null", manifest.DeletedAt, manifest.ID)
	if err != nil {
		return err
	}
	return nil
}

func GetManifestEntries(manifestId int64, orgId ...int64) (*[]ManifestEntry, error) {
	var row *sql.Rows
	var err error
	if len(orgId) > 0 {
		row, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.organization_id, e.created_at, e.updated_at from ods.manifest_entries e inner join ods.manifests m on m.id = e.manifest_id and m.deleted_at is null where e.manifest_id = ? and e.organization_id = ? and e.deleted_at is null", manifestId, orgId[0])
	} else {
		row, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.organization_id, e.created_at, e.updated_at from ods.manifest_entries e inner join ods.manifests m on m.id = e.manifest_id and m.deleted_at is null where e.manifest_id = ? and e.deleted_at is null", manifestId)
	}
	if err != nil {
		return nil, err
	}
	defer row.Close()
	entries := make([]ManifestEntry, 0)
	for row.Next() {
		manifestEntry := ManifestEntry{}
		if err = row.Scan(&manifestEntry.ID, &manifestEntry.ManifestId, &manifestEntry.AssetId, &manifestEntry.OrganizationId, &manifestEntry.CreatedAt, &manifestEntry.UpdatedAt); err != nil {
			return nil, err
		}
		entries = append(entries, manifestEntry)
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return &entries, nil
}

func GetManifestEntriesExpanded(manifestId int64, orgId ...int64) (*[]ManifestEntry, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.created_at, e.updated_at, e.organization_id, a.id, a.model_id, a.serial_number, a.room_id, a.created_at, a.updated_at, m2.id, m2.name, m2.manufacturer_id, m2.category_id, m2.created_at, m2.updated_at, c.id, c.name, c.created_at, c.updated_at, m3.id, m3.name, m3.created_at, m3.updated_at from ods.manifest_entries e inner join ods.assets a on a.id = e.asset_id inner join ods.models m2 on m2.id = a.model_id inner join ods.categories c on c.id = m2.category_id inner join ods.manufacturers m3 on m3.id = m2.manufacturer_id where e.manifest_id = ? and e.organization_id = ? and e.deleted_at is null", manifestId, orgId[0])
	} else {
		rows, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.created_at, e.updated_at, e.organization_id, a.id, a.model_id, a.serial_number, a.room_id, a.created_at, a.updated_at, m2.id, m2.name, m2.manufacturer_id, m2.category_id, m2.created_at, m2.updated_at, c.id, c.name, c.created_at, c.updated_at, m3.id, m3.name, m3.created_at, m3.updated_at from ods.manifest_entries e inner join ods.assets a on a.id = e.asset_id inner join ods.models m2 on m2.id = a.model_id inner join ods.categories c on c.id = m2.category_id inner join ods.manufacturers m3 on m3.id = m2.manufacturer_id where e.manifest_id = ? and e.deleted_at is null", manifestId)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Next()
	entries := make([]ManifestEntry, 0)
	for rows.Next() {
		entry := ManifestEntry{}
		asset := Asset{
			AssetTags: make([]AssetTag, 0),
		}
		model := Model{}
		mfg := Manufacturer{}
		cat := Category{}
		if err = rows.Scan(
			&entry.ID, &entry.ManifestId, &entry.AssetId, &entry.CreatedAt, &entry.UpdatedAt, &entry.OrganizationId, &asset.ID, &asset.ModelId, &asset.SerialNumber, &asset.RoomId, &asset.CreatedAt, &asset.UpdatedAt, &model.ID, &model.Name, &model.ManufacturerId, &model.CategoryId, &model.CreatedAt, &model.UpdatedAt, &cat.ID, &cat.Name, &cat.CreatedAt, &cat.UpdatedAt, &mfg.ID, &mfg.Name, &mfg.CreatedAt, &mfg.UpdatedAt); err != nil {
			return nil, err
		}
		model.Manufacturer = mfg
		model.Category = cat
		asset.Model = model
		assetTags, err := GetAssetTagsByAssetId(asset.ID, entry.OrganizationId)
		if err != nil {
			return nil, err
		}
		asset.AssetTags = *assetTags
		entry.Asset = asset
		entries = append(entries, entry)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &entries, nil
}

func AssetExistsInManifest(orgId int64, manifestId int64, assetId int64) (bool, error) {
	rows, err := data.GetConn().Query("select count(e.id) from ods.manifest_entries e where e.manifest_id = ? and e.organization_id = ? and e.asset_id = ? and e.deleted_at is null", manifestId, orgId, assetId)
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var c int64
		if err = rows.Scan(&c); err != nil {
			return false, err
		}
		return c > 0, nil
	}
	if err = rows.Err(); err != nil {
		return false, err
	}
	return false, errors.New("impossible state reached")
}

func AddEntryToManifest(entry *ManifestEntry) error {
	entry.ID = snowflake.GetNode().Generate().Int64()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.DeletedAt = sql.NullTime{}
	if entry.ManifestId == 0 {
		return errors.New("cannot add entry to non existent manifest")
	}
	res, err := data.GetConn().Exec("insert into ods.manifest_entries (id, manifest_id, asset_id, organization_id, created_at, updated_at) values(?,?,?,?,?,?)", entry.ID, entry.ManifestId, entry.AssetId, entry.OrganizationId, entry.CreatedAt, entry.UpdatedAt)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one row affected but got different")
	}
	return nil
}

func GetManifestEntryById(manifestId int64, id int64, orgId ...int64) (*ManifestEntry, error) {
	var row *sql.Rows
	var err error
	if len(orgId) > 0 {
		row, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.organization_id, e.created_at, e.updated_at from ods.manifest_entries e inner join ods.manifests m on m.id = e.manifest_id and m.deleted_at is null where e.id = ? and e.manifest_id = ? and e.organization_id = ? and e.deleted_at is null", id, manifestId, orgId[0])
	} else {
		row, err = data.GetConn().Query("select e.id, e.manifest_id, e.asset_id, e.organization_id, e.created_at, e.updated_at from ods.manifest_entries e inner join ods.manifests m on m.id = e.manifest_id and m.deleted_at is null where e.id = ? and e.manifest_id = ? and e.deleted_at is null", id, manifestId)
	}
	if err != nil {
		return nil, err
	}
	defer row.Close()
	if row.Next() {
		manifestEntry := ManifestEntry{}
		if err = row.Scan(&manifestEntry.ID, &manifestEntry.ManifestId, &manifestEntry.AssetId, &manifestEntry.OrganizationId, &manifestEntry.CreatedAt, &manifestEntry.UpdatedAt); err != nil {
			return nil, err
		}
		return &manifestEntry, nil
	}
	if err = row.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func DeleteManifestEntry(entry *ManifestEntry) error {
	if entry.DeletedAt.Valid {
		return errors.New("cannot delete entry that has already been deleted")
	}
	entry.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.manifest_entries e set e.deleted_at = ? where e.id = ? and e.organization_id = ? and e.deleted_at is null", entry.DeletedAt, entry.ID, entry.OrganizationId)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected one affected row but got different")
	}
	return nil
}
