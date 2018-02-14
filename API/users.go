package API

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type User struct {
	Id          int       `gorm:"column:id;not null;type:integer"`
	FirstName   string    `gorm:"column:first_name;not null;type:text"`
	SecondName  string    `gorm:"column:second_name;not null;type:text"`
	DateOfBirth time.Time `gorm:"column:date_of_birth;type:date"`
	RoleId      int       `gorm:"column:role_id;not null;type:integer"`
	Created     time.Time `gorm:"column:created;type:date"`
	Updated     time.Time `gorm:"column:updated;type:date;default:''"`
	Deleted     time.Time `gorm:"column:created;type:date;default:''"`
	Password    string    `gorm:"column:password;type:text"`
	Login       string    `gorm:"column:login;not null;type:text"`
}

func (a *API) GetUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	var users []User
	var whereString string

	id := r.FormValue("id")
	WhereBlock("id", id, &whereString)

	fName := r.FormValue("first_name")
	WhereBlock("first_name", fName, &whereString)

	sName := r.FormValue("second_name")
	WhereBlock("second_name", sName, &whereString)

	rId := r.FormValue("role_id")
	WhereBlock("role_id", rId, &whereString)

	WhereBlock("deleted", "NULL", &whereString)

	query := "SELECT * FROM users WHERE " + whereString

	err = a.Db.Raw(query).Scan(&users).Error
	if err != nil {
		a.Log.Error("Get query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Get query error | Query: " + query))
		return
	}

	json.NewEncoder(w).Encode(users)
}

func (a *API) UpdateUsers(w http.ResponseWriter, r *http.Request) {
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

	val := r.FormValue("first_name")
	SetBlock("first_name", val, &setString, true)

	val = r.FormValue("second_name")
	SetBlock("second_name", val, &setString, true)

	val = r.FormValue("date_of_birth")
	SetBlock("date_of_birth", val, &setString, true)

	val = r.FormValue("role_id")
	SetBlock("role_id", val, &setString, false)

	val = r.FormValue("password")
	passHash := sha256.Sum256([]byte(val))
	SetBlock("password", string(passHash[:]), &setString, true)

	SetBlock("updated", time.Now().Format("2006-01-02"), &setString, true)

	query := "UPDATE users SET " + setString + " WHERE " + whereString

	fmt.Println(query)

	err = a.Db.Exec(query).Error
	if err != nil {
		a.Log.Error("update query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: update query error | Query: " + query))
	}
}

func (a *API) InsertUsers(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	var tmpUsers User

	tmpUsers.FirstName = r.FormValue("first_name")
	tmpUsers.SecondName = r.FormValue("second_name")

	tmpUsers.DateOfBirth, err = time.Parse("2006-01-02", r.FormValue("date_of_birth"))
	if err != nil {
		fmt.Println("1")
		return
	}
	tmpUsers.RoleId, err = strconv.Atoi(r.FormValue("role_id"))
	if err != nil {
		return
	}

	tmpPassHash := sha256.Sum256([]byte(r.FormValue("password")))
	tmpUsers.Password = string(tmpPassHash[:])
	fmt.Println(tmpUsers.Password)
	fmt.Printf("%x", tmpUsers.Password)
	tmpUsers.Password = fmt.Sprintf("%x", tmpUsers.Password)

	tmpUsers.Created = time.Now()

	err = a.Db.Create(&tmpUsers).Error
	if err != nil {
		a.Log.Error("Insert query error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Insert query error | Query: "))
	}
	fmt.Println(tmpUsers)
}

func (a *API) DeletetUsers(w http.ResponseWriter, r *http.Request) {
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

	err = a.Db.Exec("UPDATE users SET deleted = '" + time.Now().Format("2006-01-02") + "' WHERE id in (" + id + ") AND deleted is NULL").Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: delete query error"))
	}
}
