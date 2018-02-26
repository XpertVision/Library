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

	var books []Book
	var bookTmp Book

	id := r.FormValue("id")
	if id != "" {
		bookTmp.Id, err = strconv.Atoi(id)
		if err != nil {
			a.Log.Error("problem with convert string to int (id)")
		}
	}

	bookTmp.Name = r.FormValue("name")

	bookTmp.Author = r.FormValue("author")

	uId := r.FormValue("user_id")
	if uId != "" {
		bookTmp.UserId, err = strconv.Atoi(uId)
		if err != nil {
			a.Log.Error("problem with convert string to int (user_id)")
		}
	}

	err = a.Db.Where(&bookTmp).Find(&books).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: Get query error"))
		return
	}

	json.NewEncoder(w).Encode(books)
}

func (a *API) UpdateBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	var bookTmp Book

	id := r.FormValue("id")
	bookTmp.Id, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	bookTmp.Name = r.FormValue("name")

	bookTmp.Author = r.FormValue("author")

	uId := r.FormValue("user_id")
	bookTmp.UserId, err = strconv.Atoi(uId)
	if err != nil {
		a.Log.Error("problem with convert string to int")
	}

	bookTmp.Updated = time.Now()

	err = a.Db.Model(&bookTmp).Updates(bookTmp).Error
	if err != nil {
		a.Log.Error("update query error | Query: ")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("BAD REQUEST: update query error | Query: "))
	}
}

func (a *API) InsertBooks(w http.ResponseWriter, r *http.Request) {
	var err error

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

	id := r.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: empty id"))
		return
	}

	err = a.Db.Exec("UPDATE books SET deleted = ? WHERE id = ?", time.Now().Format("2006-01-02"), id).Error
	if err != nil {
		a.Log.Error("Delete query error! Query: ")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("BAD REQUEST: delete query error"))
	}
}
