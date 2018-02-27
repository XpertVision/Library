package api

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

//User struct is struct for users table in db
type User struct {
	ID          int       `gorm:"column:id;not null;type:integer"`
	FirstName   string    `gorm:"column:first_name;not null;type:text"`
	SecondName  string    `gorm:"column:second_name;not null;type:text"`
	DateOfBirth time.Time `gorm:"column:date_of_birth;type:date"`
	RoleID      int       `gorm:"column:role_id;not null;type:integer"`
	Created     time.Time `gorm:"column:created;type:date"`
	Updated     time.Time `gorm:"column:updated;type:date;default:''"`
	Deleted     time.Time `gorm:"column:created;type:date;default:''"`
	Password    string    `gorm:"column:password;type:text"`
	Login       string    `gorm:"column:login;not null;type:text"`
}

//GetUsers func return to web users table data
func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	var users []User
	var userTmp User

	id := r.FormValue("id")
	userTmp.ID, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int (id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	userTmp.FirstName = r.FormValue("first_name")

	userTmp.SecondName = r.FormValue("second_name")

	rID := r.FormValue("role_id")
	userTmp.RoleID, err = strconv.Atoi(rID)
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	err = a.DB.Where(&userTmp).Find(&users).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	json.NewEncoder(w).Encode(users)
}

//UpdateUsers func updates data in users table
func (a *API) UpdateUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	var userTmp User
	id := r.FormValue("id")
	userTmp.ID, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int (id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	userTmp.FirstName = r.FormValue("first_name")

	userTmp.SecondName = r.FormValue("second_name")

	dofBTmp := r.FormValue("date_of_birth")
	userTmp.DateOfBirth, err = time.Parse("2006-01-02", dofBTmp)
	if err != nil {
		a.Log.Error("problem with convert string to time.time (date_of_birth) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	rID := r.FormValue("role_id")
	userTmp.RoleID, err = strconv.Atoi(rID)
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	userTmp.Login = r.FormValue("login")

	passFromForm := r.FormValue("password")
	if passFromForm != "" {
		tmpPassHash := sha256.Sum256([]byte(r.FormValue("password")))
		tmpPasswordStr := string(tmpPassHash[:])
		tmpPasword := fmt.Sprintf("%x", tmpPasswordStr)
		userTmp.Password = tmpPasword
	}

	userTmp.Updated = time.Now()

	err = a.DB.Model(&userTmp).Updates(userTmp).Error
	if err != nil {
		a.Log.Error("update query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//InsertUsers func inserts data in users table
func (a *API) InsertUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpUsers User

	tmpUsers.FirstName = r.FormValue("first_name")
	tmpUsers.SecondName = r.FormValue("second_name")

	tmpUsers.DateOfBirth, err = time.Parse("2006-01-02", r.FormValue("date_of_birth"))
	if err != nil {
		a.Log.Error("problem with convert string to time.time (date_of_birth) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	tmpUsers.RoleID, err = strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		a.Log.Error("problem with convert string to int (role_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	tmpUsers.Login = r.FormValue("login")

	tmpPassHash := sha256.Sum256([]byte(r.FormValue("password")))
	tmpUsers.Password = string(tmpPassHash[:])
	tmpUsers.Password = fmt.Sprintf("%x", tmpUsers.Password)

	tmpUsers.Created = time.Now()

	err = a.DB.Create(&tmpUsers).Error
	if err != nil {
		a.Log.Error("Insert query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//DeleteUsers func set delete column in users table for row
func (a *API) DeleteUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	id := r.FormValue("id")
	if id == "" {
		a.Log.Error("Empty id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}

	err = a.DB.Exec("UPDATE users SET deleted = ? WHERE id = ? AND deleted IS NULL", time.Now().Format("2006-01-02"), id).Error
	if err != nil {
		a.Log.Error("Delete query error | Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}
}
