package api

import (
	"bytes"
	"encoding/json"
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var testApi API

func TestGetBooks(t *testing.T) {
	var err error

	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testApi.Db = db
	testApi.Log = loger

	var bookArrayFirst, bookArrayLast []Book
	db.Raw("SELECT * FROM books WHERE deleted IS NULL").Scan(&bookArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getBooks", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.GetBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(http.StatusOK, w.Code)
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

	loger, err := lumber.NewRotateLogger("log_"+time.Now().Format("2006-01-02")+".log", 10000, 10)
	if err != nil {
		t.Fatal(err)
	}

	db, err := gorm.Open("postgres", "host=localhost port=5433 user=admin dbname=GOLibrary sslmode=disable password=12344321Qw5")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testApi.Db = db
	testApi.Log = loger

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/InsertBooks?name=TEST_FOR_TEST&author=TESTER&user_id=777", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertBooks(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(http.StatusOK, w.Code)
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
