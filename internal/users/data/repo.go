package data

import (
	"convention.ninja/internal/data"
	"convention.ninja/internal/snowflake"
	"database/sql"
	"errors"
	"time"
)

// TODO implement get user by id

func GetUserByFirebase(firebaseId string) (*User, error) {
	rows, err := data.GetConn().Query("select id,name,display_name,email,email_verified,firebase_id,created_at,updated_at,deleted_at from ods.users where firebase_id = ? and deleted_at is null", firebaseId)
	if err != nil {
		return nil, err
	}
	var user User
	if rows.Next() {
		if err = rows.Scan(&user.ID, &user.Name, &user.DisplayName, &user.Email, &user.EmailVerified, &user.FirebaseId, &user.CreatedAt, &user.UpdatedAt, &user.DeletedAt); err != nil {
			return nil, err
		}
		return &user, nil
	}
	return nil, nil
}

func EmailExists(email string) (bool, error) {
	rows, err := data.GetConn().Query("select count(*) from ods.users where email = ? and deleted_at is null", email)
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
	return false, errors.New("got no records but expected one")
}

func CreateUser(user *User) error {
	user.ID = snowflake.GetNode().Generate().Int64()
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.DeletedAt = sql.NullTime{}

	res, err := data.GetConn().Exec("insert into ods.users (id,name,display_name,email,email_verified,firebase_id,created_at,updated_at) VALUES(?,?,?,?,?,?,?,?)", user.ID, user.Name, user.DisplayName, user.Email, user.EmailVerified, user.FirebaseId, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil || affected == 0 {
		return errors.New("expected 1 affected row got 0")
	}
	return nil
}

func UpdateUser(user *User) error {
	user.UpdatedAt = time.Now()
	if user.DeletedAt.Valid {
		return errors.New("cannot update user who is already deleted")
	}
	res, err := data.GetConn().Exec("update ods.users set name = ?, display_name = ?, email = ?, email_verified = ?, updated_at = ? where id = ? and deleted_at is null limit 1", user.Name, user.DisplayName, user.Email, user.EmailVerified, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("expected 1 affected row but got different result")
	}
	return nil
}

func DeleteUser(user *User) error {
	if user.DeletedAt.Valid {
		return errors.New("cannot delete user who is already deleted")
	}
	user.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	res, err := data.GetConn().Exec("update ods.users set deleted_at = ? where id = ? and deleted_at is null limit 1", user.DeletedAt, user.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil || count != 1 {
		return errors.New("expected 1 affected row but got different result")
	}
	return nil
}
