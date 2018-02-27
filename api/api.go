package api

import (
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
)

//API struct than contain DB and Loger objects
type API struct {
	DB  *gorm.DB
	Log lumber.Logger
}

//New return new API object
func New(db *gorm.DB, log lumber.Logger) *API {
	api := new(API)

	api.DB = db
	api.Log = log

	return api
}
