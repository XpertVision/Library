package api

import (
	"errors"
	"time"
)

//Connection struct is struct for connections table in db
type Connection struct {
	ID           int       `gorm:"column:id;type:integer"`
	UserID       int       `gorm:"column:user_id;type:integer"`
	Token        string    `gorm:"column:token;type:text"`
	RoleID       int       `gorm:"column:role_id;type:integer"`
	GenerateDate time.Time `gorm:"column:generate_date;type:timestamp without time zone"`
}

//GetConnectionFromToken func return connection struct with full data from connections table with "token" filter
func (a *API) GetConnectionFromToken(token string) (Connection, error) {
	var err error
	var conn Connection

	err = a.DB.Find(&conn).Where("token = ?", token).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		return conn, errors.New("Get connection from db error")
	}

	return conn, nil
}

//GetConnectionFromID func return connection struct with full data from connections table with "userID" filter
func (a *API) GetConnectionFromID(userID int) (Connection, error) {
	var err error
	var conn Connection

	err = a.DB.Find(&conn).Where("user_id = ?", userID).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		return conn, errors.New("Get connection from db error")
	}

	return conn, nil
}

//InsertConnection func inserts data in connections table
func (a *API) InsertConnection(conn Connection) error {
	var err error

	err = a.DB.Create(&conn).Error
	if err != nil {
		a.Log.Error("Insert query error | Error: ", err)
		return errors.New("Insert connection to db error")
	}

	return nil
}

//UpdateConnection func updates data in connections table
func (a *API) UpdateConnection(conn Connection) error {
	var err error

	err = a.DB.Model(&conn).Where("id = ?", conn.ID).Update(&conn).Error
	if err != nil {
		a.Log.Error("update query error | Error: ", err)
		return errors.New("Update connection in db error")
	}

	return nil
}

//DeleteConnection func delete rom from connections table where target token
func (a *API) DeleteConnection(token string) error {
	var err error

	err = a.DB.Exec("DELETE FROM connections WHERE token = ?", token).Error
	if err != nil {
		a.Log.Error("Delete query error! Error: ", err)
		return errors.New("Delete connection from db error")
	}

	return nil
}
