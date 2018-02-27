package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

//Book struct is struct for books table in db
type Book struct {
	ID      int       `gorm:"column:id;not null;type:integer"`
	Name    string    `gorm:"column:name;not null;type:text"`
	Author  string    `gorm:"column:author;not null;type:text"`
	UserID  int       `gorm:"column:user_id;not null;type:integer"`
	Created time.Time `gorm:"column:created;type:date"`
	Updated time.Time `gorm:"column:updated;type:date;default:''"`
	Deleted time.Time `gorm:"column:deleted;type:date;default:''"`
}

//GetBooks func return to web books table data
func (a *API) GetBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	var books []Book
	var bookTmp Book

	id := r.FormValue("id")
	if id != "" {
		bookTmp.ID, err = strconv.Atoi(id)
		if err != nil {
			a.Log.Error("problem with convert string to int (id) | Error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR"))
			return
		}
	}

	bookTmp.Name = r.FormValue("name")

	bookTmp.Author = r.FormValue("author")

	uID := r.FormValue("user_id")
	if uID != "" {
		bookTmp.UserID, err = strconv.Atoi(uID)
		if err != nil {
			a.Log.Error("problem with convert string to int (user_id) | Error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR"))
			return
		}
	}

	err = a.DB.Where(&bookTmp).Find(&books).Error
	if err != nil {
		a.Log.Error("Get query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	json.NewEncoder(w).Encode(books)
}

//UpdateBooks func updates data in books table
func (a *API) UpdateBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	var bookTmp Book

	id := r.FormValue("id")
	bookTmp.ID, err = strconv.Atoi(id)
	if err != nil {
		a.Log.Error("problem with convert string to int (id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	bookTmp.Name = r.FormValue("name")

	bookTmp.Author = r.FormValue("author")

	uID := r.FormValue("user_id")
	bookTmp.UserID, err = strconv.Atoi(uID)
	if err != nil {
		a.Log.Error("problem with convert string to int (user_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	bookTmp.Updated = time.Now()

	err = a.DB.Model(&bookTmp).Updates(bookTmp).Error
	if err != nil {
		a.Log.Error("update query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//InsertBooks func inserts data in books table
func (a *API) InsertBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	var tmpBooks Book

	tmpBooks.Name = r.FormValue("name")
	tmpBooks.Author = r.FormValue("author")
	tmpBooks.UserID, err = strconv.Atoi(r.FormValue("user_id"))
	if err != nil {
		a.Log.Error("problem with convert string to int (user_id) | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
	tmpBooks.Created = time.Now()

	err = a.DB.Create(&tmpBooks).Error
	if err != nil {
		a.Log.Error("Insert query error | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}
}

//DeleteBooks func set delete column in books table for row
func (a *API) DeleteBooks(w http.ResponseWriter, r *http.Request) {
	var err error

	id := r.FormValue("id")
	if id == "" {
		a.Log.Error("empty id")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}

	err = a.DB.Exec("UPDATE books SET deleted = ? WHERE id = ?", time.Now().Format("2006-01-02"), id).Error
	if err != nil {
		a.Log.Error("Delete query error! Error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("ERROR"))
		return
	}
}
