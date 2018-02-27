package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
)

var testAPI API

func initAPI(a *API, t *testing.T) {
	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		t.Fatal(err)
	}

	a.DB = db
	a.Log = loger
}

func TestGetBooks(t *testing.T) {
	var err error
	var bookArrayFirst, bookArrayLast []Book

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	testAPI.DB.Raw("SELECT * FROM books WHERE deleted IS NULL").Scan(&bookArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getBooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.GetBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	firstBytes, err := json.Marshal(bookArrayFirst)
	if err != nil {
		t.Error(err)
	}

	json.NewDecoder(w.Body).Decode(&bookArrayLast)
	lastBytes, err := json.Marshal(bookArrayLast)
	if err != nil {
		t.Error(err)
	}

	if !bytes.Equal(firstBytes, lastBytes) {
		t.Error("error, not equal")
	}
}

func TestInsertBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testAPI.DB.Raw("SELECT * FROM books WHERE name = 'TEST_FOR_TEST' AND author = 'TESTER' AND user_id = 777 AND deleted IS NULL").Scan(&booksFromDb)

	if cap(booksFromDb) != 1 {
		t.Error("too much rows found")
	}

	if booksFromDb[0].UserID != 777 {
		t.Error("Wrong user_id")
	}

	if booksFromDb[0].Name != "TEST_FOR_TEST" {
		t.Error("Wrong name")
	}

	if booksFromDb[0].Author != "TESTER" {
		t.Error("Wrong author")
	}

	if booksFromDb[0].Created.IsZero() {
		t.Error("Wrong create time")
	}

	testAPI.DB.Exec("DELETE FROM books WHERE name = 'TEST_FOR_TEST' AND author = 'TESTER' AND user_id = 777")
}

func TestDeleteBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testAPI.DB.Raw("SELECT * FROM books WHERE user_id=777 AND deleted IS NULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/deleteBooks?id="+strconv.Itoa(booksFromDb[0].ID)+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.DeleteBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testAPI.DB.Raw("SELECT * FROM books WHERE user_id=777 AND deleted NOTNULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't delete")
	}

	testAPI.DB.Exec("DELETE FROM books WHERE user_id=777 AND deleted NOTNULL")
}

func TestUpdateBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testAPI.DB.Raw("SELECT * FROM books WHERE user_id=777 AND deleted IS NULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/updateBooks?id="+strconv.Itoa(booksFromDb[0].ID)+"&name=TESTN&author=TESTA&user_id=888", nil)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.UpdateBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testAPI.DB.Raw("SELECT * FROM books WHERE id=" + strconv.Itoa(booksFromDb[0].ID) + " AND updated NOTNULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't update")
	}

	if booksFromDb[0].Name != "TESTN" {
		t.Error("Wrong updated name")
	}

	if booksFromDb[0].Author != "TESTA" {
		t.Error("Wrong updated author")
	}

	if booksFromDb[0].UserID != 888 {
		t.Error("Wrong updated user_id")
	}

	testAPI.DB.Exec("DELETE FROM books WHERE id = " + strconv.Itoa(booksFromDb[0].ID) + " AND updated NOTNULL")
}
