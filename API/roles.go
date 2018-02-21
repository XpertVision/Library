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

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	var roles []Role
	var whereString string

	id := r.FormValue("id")
	WhereBlock("id", id, &whereString)

	name := r.FormValue("name")
	WhereBlock("name", name, &whereString)

	roleId := r.FormValue("role_id")
	WhereBlock("role_id", roleId, &whereString)

	WhereBlock("deleted", "NULL", &whereString)

	query := "SELECT * FROM roles WHERE " + whereString

	err = a.Db.Raw(query).Scan(&roles).Error
	if err != nil {
		a.Log.Error("Get query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Get query error | Query: " + query))
		return
	}

	json.NewEncoder(w).Encode(roles)
}

func (a *API) UpdateRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: empty id"))
		return
	}

	var whereString string
	var setString string

	WhereBlock("id", id, &whereString)
	WhereBlock("deleted", "NULL", &whereString)

	val := r.FormValue("name")
	SetBlock("name", val, &setString, true)

	val = r.FormValue("role_id")
	SetBlock("role_id", val, &setString, false)

	val = r.FormValue("allowe_paths")
	SetBlock("allowe_paths", val, &setString, true)

	SetBlock("updated", time.Now().Format("2006-01-02"), &setString, true)

	query := "UPDATE roles SET " + setString + " WHERE " + whereString

	err = a.Db.Exec(query).Error
	if err != nil {
		a.Log.Error("update query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: update query error | Query: " + query))
	}
}

func (a *API) InsertRoles(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

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

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	id := r.FormValue("id")

	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: empty id"))
		return
	}

	err = a.Db.Exec("UPDATE roles SET deleted = '" + time.Now().Format("2006-01-02") + "' WHERE id in (" + id + ") AND deleted is NULL").Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: delete query error"))
	}
}

func (a *API) GetRoleFromRoleId(roleId int) (Role, error) {
	var err error
	var roleTmp Role

	err = a.Db.Raw("SELECT * FROM roles WHERE role_id = " + strconv.Itoa(roleId)).Scan(&roleTmp).Error
	if err != nil {
		fmt.Println("24")
		a.Log.Error("Delete query error! Query: ")
		return roleTmp, errors.New("Delete connection from db error")
	}

	return roleTmp, nil
}
