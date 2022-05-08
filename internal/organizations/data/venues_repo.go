package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

func GetVenuesByOrganization(orgId int64) (*[]Venue, error) {
	rows, err := data.GetConn().Query("select id, name, address, organization_id, created_at, updated_at, deleted_at from ods.venues where organization_id = ? and deleted_at is null", orgId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []Venue
	results = make([]Venue, 0)
	for rows.Next() {
		venue := Venue{}
		venue.Rooms = make([]VenueRoom, 0)
		if err = rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.OrganizationId, &venue.CreatedAt, &venue.UpdatedAt, &venue.DeletedAt); err != nil {
			return nil, err
		}
		results = append(results, venue)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &results, nil
}

func GetVenueById(id int64, orgId ...int64) (*Venue, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, address, organization_id, created_at, updated_at, deleted_at from ods.venues where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select id, name, address, organization_id, created_at, updated_at, deleted_at from ods.venues where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		venue := Venue{}
		if err = rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.OrganizationId, &venue.CreatedAt, &venue.UpdatedAt, &venue.DeletedAt); err != nil {
			return nil, err
		}
		return &venue, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func CreateVenue(venue *Venue) error {
	venue.ID = snowflake.GetNode().Generate().Int64()
	venue.CreatedAt = time.Now()
	venue.UpdatedAt = time.Now()
	venue.DeletedAt = sql.NullTime{}
	res, err := data.GetConn().Exec("insert into ods.venues (id, name, address, organization_id, created_at, updated_at) values (?, ?, ?, ?, ?, ?)", venue.ID, venue.Name, venue.Address, venue.OrganizationId, venue.CreatedAt, venue.UpdatedAt)
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

func UpdateVenue(venue *Venue) error {
	if venue.DeletedAt.Valid {
		return errors.New("cannot update venue that has already been deleted")
	}
	venue.UpdatedAt = time.Now()
	res, err := data.GetConn().Exec("update ods.venues set name = ?, address = ? where id = ? and organization_id = ? and deleted_at is null limit 1", venue.Name, venue.Address, venue.ID, venue.OrganizationId)
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

func DeleteVenue(venue *Venue) error {
	if venue.DeletedAt.Valid {
		return errors.New("cannot delete venue that has already been deleted")
	}
	venue.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.venues set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", venue.DeletedAt, venue.ID, venue.OrganizationId)
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

func GetVenueRooms(venueId int64, orgId ...int64) (*[]VenueRoom, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, venue_id, organization_id, created_at, updated_at, deleted_at from ods.venue_rooms where organization_id = ? and venue_id = ? and deleted_at is null", orgId[0], venueId)
	} else {
		rows, err = data.GetConn().Query("select id, name, venue_id, organization_id, created_at, updated_at, deleted_at from ods.venue_rooms where venue_id = ? and deleted_at is null", venueId)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results = make([]VenueRoom, 0)
	for rows.Next() {
		room := VenueRoom{}
		if err = rows.Scan(&room.ID, &room.Name, &room.VenueId, &room.OrganizationId, &room.CreatedAt, &room.UpdatedAt, &room.DeletedAt); err != nil {
			return nil, err
		}
		results = append(results, room)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &results, nil
}

func GetVenueRoomById(id int64, orgId ...int64) (*VenueRoom, error) {
	var rows *sql.Rows
	var err error
	if len(orgId) > 0 {
		rows, err = data.GetConn().Query("select id, name, venue_id, organization_id, created_at, updated_at, deleted_at from ods.venue_rooms where organization_id = ? and id = ? and deleted_at is null", orgId[0], id)
	} else {
		rows, err = data.GetConn().Query("select id, name, venue_id, organization_id, created_at, updated_at, deleted_at from ods.venue_rooms where id = ? and deleted_at is null", id)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	if rows.Next() {
		room := VenueRoom{}
		if err = rows.Scan(&room.ID, &room.Name, &room.VenueId, &room.OrganizationId, &room.CreatedAt, &room.UpdatedAt, &room.DeletedAt); err != nil {
			return nil, err
		}
		return &room, nil
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}

func CreateVenueRoom(room *VenueRoom) error {
	room.ID = snowflake.GetNode().Generate().Int64()
	room.CreatedAt = time.Now()
	room.UpdatedAt = time.Now()
	room.DeletedAt = sql.NullTime{}
	res, err := data.GetConn().Exec("insert into ods.venue_rooms (id, name, venue_id, organization_id, created_at, updated_at) values(?, ?, ?, ?, ?, ?)", room.ID, room.Name, room.VenueId, room.OrganizationId, room.CreatedAt, room.UpdatedAt)
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

func DeleteVenueRoom(room *VenueRoom) error {
	if room.DeletedAt.Valid {
		return errors.New("cannot delete room that has already been deleted")
	}
	room.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.venue_rooms set deleted_at = ? where id = ? and organization_id = ? and deleted_at is null limit 1", room.DeletedAt, room.ID, room.OrganizationId)
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
