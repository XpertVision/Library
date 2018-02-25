package api

import (
	"errors"
	"time"
)

type Connection struct {
	Id           int       `gorm:"column:id;type:integer"`
	UserId       int       `gorm:"column:user_id;type:integer"`
	Token        string    `gorm:"column:token;type:text"`
	RoleId       int       `gorm:"column:role_id;type:integer"`
	GenerateDate time.Time `gorm:"column:generate_date;type:timestamp without time zone"`
}

func (a *API) GetConnectionFromToken(token string) (Connection, error) {
	var err error
	var conn Connection

	err = a.Db.Find(&conn).Where("token = ?", token).Error
	if err != nil {
		a.Log.Error("Get query error | Query: " /* + query*/)
		return conn, errors.New("Get connection from db error")
	}

	return conn, nil
}

func (a *API) GetConnectionFromId(userId int) (Connection, error) {
	var err error
	var conn Connection

	err = a.Db.Find(&conn).Where("user_id = ?", userId).Error
	if err != nil {
		a.Log.Error("Get query error | Query: " /* + query*/)
		return conn, errors.New("Get connection from db error")
	}

	return conn, nil
}

func (a *API) InsertConnection(conn Connection) error {
	var err error

	err = a.Db.Create(&conn).Error
	if err != nil {
		a.Log.Error("Insert query error")
		return errors.New("Insert connection to db error")
	}

	return nil
}

func (a *API) UpdateConnection(conn Connection) error {
	var err error

	err = a.Db.Model(&conn).Where("id = ?", conn.Id).Update(&conn).Error
	if err != nil {
		a.Log.Error("update query error | Query: " /* + query*/)
		return errors.New("Update connection in db error")
	}

	return nil
}

func (a *API) DeleteConnection(token string) error {
	var err error

	err = a.Db.Exec("DELETE FROM connections WHERE token = ?", token).Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		return errors.New("Delete connection from db error")
	}

	return nil
}
