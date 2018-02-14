package API

import (
	"encoding/json"
	"net/http"
	"time"
)

type Role struct {
	Id   int    `gorm:"column:id;not null;type:integer"`
	Name string `gorm:"column:name;not null;type:text"`
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

	WhereBlock("deleted", "NULL", &whereString)

	query := "SELECT * FROM books WHERE " + whereString

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

	query := "UPDATE books SET " + setString + " WHERE " + whereString

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
