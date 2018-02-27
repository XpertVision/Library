package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/lib/pq"
)

//Role struct is struct for roles table in db
type Role struct {
	ID          int            `gorm:"column:id;not null;type:integer"`
	Name        string         `gorm:"column:name;not null;type:text"`
	RoleID      int            `gorm:"column:role_id;not null;type:integer"`
	AllowePaths pq.StringArray `gorm:"column:allowe_paths;type:text[]"`
	Deleted     time.Time      `gorm:"column:deleted;type:date;default:''"`
	Updated     time.Time      `gorm:"column:updated;type:date;default:''"`
	Created     time.Time      `gorm:"column:created;type:date;default:''"`
}

//GetRoles func return to web roles table data
func (a *API) GetRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var roles []Role
	var roleTmp Role

	id := r.FormValue("id")
	roleTmp.ID, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int (id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	roleTmp.Name = r.FormValue("name")

	roleID := r.FormValue("role_id")
	roleTmp.ID, err = strconv.Atoi(roleID)
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	err = a.DB.Where(&roleTmp).Find(&roles).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	json.NewEncoder(w).Encode(roles)
}

//UpdateRoles func updates data in roles table
func (a *API) UpdateRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var roleTmp Role

	id := r.FormValue("id")
	roleTmp.ID, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int (id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	roleTmp.Name = r.FormValue("name")

	roleID := r.FormValue("role_id")
	roleTmp.RoleID, err = strconv.Atoi(roleID)
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	allowePathTmp := r.FormValue("allowe_paths")
	roleTmp.AllowePaths.Scan(allowePathTmp)

	roleTmp.Updated = time.Now()

	err = a.DB.Model(&roleTmp).Updates(roleTmp).Error
	if err != nil {
		a.Log.Error("update query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//InsertRoles func inserts data in roles table
func (a *API) InsertRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpRoles Role

	tmpRoles.Name = r.FormValue("name")

	tmpAllowePaths := r.FormValue("allowe_paths")
	tmpRoles.AllowePaths.Scan(tmpAllowePaths)

	roleID, err := strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
	tmpRoles.RoleID = roleID

	tmpRoles.Created = time.Now()

	err = a.DB.Create(&tmpRoles).Error
	if err != nil {
		a.Log.Error("Insert query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//DeleteRoles func set delete column in roles table for row
func (a *API) DeleteRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	id := r.FormValue("id")
	if id == "" {
		a.Log.Error("empty id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}

	err = a.DB.Exec("UPDATE roles SET deleted = ? WHERE id = ? AND deleted IS NULL", time.Now().Format("2006-01-02"), id).Error
	if err != nil {
		a.Log.Error("Delete query error! Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}
}

//GetRoleFromRoleID func return role struct with full data from roles table with "roleID" filter
func (a *API) GetRoleFromRoleID(roleID int) (Role, error) {
	var err error
	var roleTmp Role

	err = a.DB.Find(&roleTmp).Where("role_id = ?", roleID).Error
	if err != nil {
		a.Log.Error("Delete query error! Error: ", err)
		return roleTmp, errors.New("Delete connection from db error")
	}

	return roleTmp, nil
}
