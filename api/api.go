package api

import (
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
)

type API struct {
	Db  *gorm.DB
	Log lumber.Logger
}

func (a *API) New(db *gorm.DB, log lumber.Logger) {
	a.Db = db
	a.Log = log
}
