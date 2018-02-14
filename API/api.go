package API

import (
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
)

type API struct {
	Db  *gorm.DB
	Log lumber.Logger
}
