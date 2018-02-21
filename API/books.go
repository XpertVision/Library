package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Book struct {
	Id      int       `gorm:"column:id;not null;type:integer"`
	Name    string    `gorm:"column:name;not null;type:text"`
	Author  string    `gorm:"column:author;not null;type:text"`
	UserId  int       `gorm:"column:user_id;not null;type:integer"`
	Created time.Time `gorm:"column:created;type:date"`
	Updated time.Time `gorm:"column:updated;type:date;default:''"`
	Deleted time.Time `gorm:"column:deleted;type:date;default:''"`
}

func (a *API) GetBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	var books []Book
	var whereString string

	id := r.FormValue("id")
	WhereBlock("id", id, &whereString)

	name := r.FormValue("name")
	WhereBlock("name", name, &whereString)

	author := r.FormValue("author")
	WhereBlock("author", author, &whereString)

	userId := r.FormValue("user_id")
	WhereBlock("user_id", userId, &whereString)

	WhereBlock("deleted", "NULL", &whereString)

	query := "SELECT * FROM books WHERE " + whereString

	err = a.Db.Raw(query).Scan(&books).Error
	if err != nil {
		a.Log.Error("Get query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Get query error | Query: " + query))
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (a *API) UpdateBooks(w http.ResponseWriter, r *http.Request) {
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

	val = r.FormValue("author")
	SetBlock("author", val, &setString, true)

	val = r.FormValue("user_id")
	SetBlock("user_id", val, &setString, false)

	SetBlock("updated", time.Now().Format("2006-01-02"), &setString, true)

	query := "UPDATE books SET " + setString + " WHERE " + whereString
	err = a.Db.Exec(query).Error
	if err != nil {
		a.Log.Error("update query error | Query: " + query)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: update query error | Query: " + query))
	}
}

func (a *API) InsertBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		a.Log.Error("Parse form error")
		return
	}

	var tmpBooks Book

	tmpBooks.Name = r.FormValue("name")
	tmpBooks.Author = r.FormValue("author")
	tmpBooks.UserId, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		a.Log.Error("Insert query error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: parse iser_id error"))
		return
	}
	tmpBooks.Created = time.Now()

	err = a.Db.Create(&tmpBooks).Error
	if err != nil {
		a.Log.Error("Insert query error")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Insert query error | Query: "))
	}
}
func (a *API) DeleteBooks(w http.ResponseWriter, r *http.Request) {
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

	err = a.Db.Exec("UPDATE books SET deleted = '" + time.Now().Format("2006-01-02") + "' WHERE id in (" + id + ") AND deleted is NULL").Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: delete query error"))
	}
}
