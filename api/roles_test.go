package api

import (
	"bytes"
	"encoding/json"
	"github.com/jcelliott/lumber"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestGetRoles(t *testing.T) {
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

	var roleArrayFirst, roleArrayLast []Role
	db.Raw("SELECT * FROM roles WHERE deleted IS NULL").Scan(&roleArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getRoles", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.GetRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(http.StatusOK, w.Code)
	}

	firstBytes, err := json.Marshal(roleArrayFirst)
	if err != nil {
		t.Error(err)
	}

	json.NewDecoder(w.Body).Decode(&roleArrayLast)
	lastBytes, err := json.Marshal(roleArrayLast)

	if !bytes.Equal(firstBytes, lastBytes) {
		t.Error("error, not equal")
	}
}

func TestInsertRoles(t *testing.T) {
	var err error
	var roleFromDb []Role
	var tmpAllowePaths pq.StringArray

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
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertRoles?name=TEST&role_id=777&allowe_paths={/TEST1,/TEST2,/TEST3}", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(http.StatusOK, w.Code)
	}

	testApi.Db.Raw("SELECT * FROM roles WHERE name = 'TEST' AND role_id = 777 AND deleted IS NULL").Scan(&roleFromDb)

	if cap(roleFromDb) != 1 {
		t.Error("too much rows found")
	}

	if roleFromDb[0].RoleId != 777 {
		t.Error("Wrong role_id")
	}

	if roleFromDb[0].Name != "TEST" {
		t.Error("Wrong name")
	}

	tmpAllowePathsStr := "{/TEST1,/TEST2,/TEST3}"
	tmpAllowePaths.Scan(tmpAllowePathsStr)

	if !reflect.DeepEqual(roleFromDb[0].AllowePaths, tmpAllowePaths) {
		t.Error("Wrong allowe_paths")
	}

	testApi.Db.Exec("DELETE FROM roles WHERE name = 'TEST' AND role_id = 777")
}
