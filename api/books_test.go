package api

import (
	"bytes"
	"encoding/json"
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"
)

var testApi API

func initApi(a *API, t *testing.T) {
	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		t.Fatal(err)
	}

	a.Db = db
	a.Log = loger
}

func TestGetBooks(t *testing.T) {
	var err error
	var bookArrayFirst, bookArrayLast []Book

	initApi(&testApi, t)
	defer testApi.Db.Close()

	testApi.Db.Raw("SELECT * FROM books WHERE deleted IS NULL").Scan(&bookArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getBooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.GetBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	firstBytes, err := json.Marshal(bookArrayFirst)
	if err != nil {
		t.Error(err)
	}

	json.NewDecoder(w.Body).Decode(&bookArrayLast)
	lastBytes, err := json.Marshal(bookArrayLast)

	if !bytes.Equal(firstBytes, lastBytes) {
		t.Error("error, not equal")
	}
}

func TestInsertBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM books WHERE name = 'TEST_FOR_TEST' AND author = 'TESTER' AND user_id = 777 AND deleted IS NULL").Scan(&booksFromDb)

	if cap(booksFromDb) != 1 {
		t.Error("too much rows found")
	}

	if booksFromDb[0].UserId != 777 {
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

	testApi.Db.Exec("DELETE FROM books WHERE name = 'TEST_FOR_TEST' AND author = 'TESTER' AND user_id = 777")
}

func TestDeleteBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM books WHERE user_id=777 AND deleted IS NULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/deleteBooks?id="+strconv.Itoa(booksFromDb[0].Id)+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.DeleteBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM books WHERE user_id=777 AND deleted NOTNULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't delete")
	}

	testApi.Db.Exec("DELETE FROM books WHERE user_id=777 AND deleted NOTNULL")
}

func TestUpdateBooks(t *testing.T) {
	var err error
	var booksFromDb []Book

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM books WHERE user_id=777 AND deleted IS NULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/updateBooks?id="+strconv.Itoa(booksFromDb[0].Id)+"&name=TESTN&author=TESTA&user_id=888", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.UpdateBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM books WHERE id=" + strconv.Itoa(booksFromDb[0].Id) + " AND updated NOTNULL").Scan(&booksFromDb)
	if cap(booksFromDb) != 1 {
		t.Error("row didn't update")
	}

	if booksFromDb[0].Name != "TESTN" {
		t.Error("Wrong updated name")
	}

	if booksFromDb[0].Author != "TESTA" {
		t.Error("Wrong updated author")
	}

	if booksFromDb[0].UserId != 888 {
		t.Error("Wrong updated user_id")
	}

	testApi.Db.Exec("DELETE FROM books WHERE id = " + strconv.Itoa(booksFromDb[0].Id) + " AND updated NOTNULL")
}
