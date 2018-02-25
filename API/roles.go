package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"net/http"
	"strconv"
	"time"
)

type Role struct {
	Id          int            `gorm:"column:id;not null;type:integer"`
	Name        string         `gorm:"column:name;not null;type:text"`
	RoleId      int            `gorm:"column:role_id;not null;type:integer"`
	AllowePaths pq.StringArray `gorm:"column:allowe_paths;type:text[]"`
	Deleted     time.Time      `gorm:"column:deleted;type:date;default:''"`
	Updated     time.Time      `gorm:"column:updated;type:date;default:''"`
	Created     time.Time      `gorm:"column:created;type:date;default:''"`
}

func (a *API) GetRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var roles []Role
	var roleTmp Role

	id := r.FormValue("id")
	roleTmp.Id, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	roleTmp.Name = r.FormValue("name")

	roleId := r.FormValue("role_id")
	roleTmp.Id, err = strconv.Atoi(roleId)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	err = a.Db.Where(&roleTmp).Find(&roles).Error
	if err != nil {
		a.Log.Error("Get query error | Query: ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Get query error"))
		return
	}

	json.NewEncoder(w).Encode(roles)
}

func (a *API) UpdateRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var roleTmp Role

	id := r.FormValue("id")
	roleTmp.Id, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	roleTmp.Name = r.FormValue("name")

	roleId := r.FormValue("role_id")
	roleTmp.RoleId, err = strconv.Atoi(roleId)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	allowePathTmp := r.FormValue("allowe_paths")
	roleTmp.AllowePaths.Scan(allowePathTmp)

	roleTmp.Updated = time.Now()

	err = a.Db.Model(&roleTmp).Updates(roleTmp).Error
	if err != nil {
		a.Log.Error("update query error | Query: ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: update query error"))
	}
}

func (a *API) InsertRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpRoles Role

	tmpRoles.Name = r.FormValue("name")

	tmpAllowePaths := r.FormValue("allowe_paths")
	tmpRoles.AllowePaths.Scan(tmpAllowePaths)

	role, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		a.Log.Error("Convert role_id error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal problems"))
		return
	}
	tmpRoles.RoleId = role

	tmpRoles.Created = time.Now()

	err = a.Db.Create(&tmpRoles).Error
	if err != nil {
		a.Log.Error("Insert query error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Insert query error | Query: "))
	}
}

func (a *API) DeleteRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: empty id"))
		return
	}

	err = a.Db.Exec("UPDATE roles SET deleted = ? WHERE id = ? AND deleted IS NULL", time.Now().Format("2006-01-02"), id).Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: delete query error"))
	}
}

func (a *API) GetRoleFromRoleId(roleId int) (Role, error) {
	var err error
	var roleTmp Role

	err = a.Db.Find(&roleTmp).Where("role_id = ?", roleId).Error
	if err != nil {
		fmt.Println("24")
		a.Log.Error("Delete query error! Query: ")
		return roleTmp, errors.New("Delete connection from db error")
	}

	return roleTmp, nil
}
