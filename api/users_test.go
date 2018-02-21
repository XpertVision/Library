package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

func TestGetUsers(t *testing.T) {
	var err error
	var userArrayFirst, userArrayLast []User

	initApi(&testApi, t)
	defer testApi.Db.Close()

	testApi.Db.Raw("SELECT * FROM users WHERE deleted IS NULL").Scan(&userArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getUsers", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.GetUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	firstBytes, err := json.Marshal(userArrayFirst)
	if err != nil {
		t.Error(err)
	}

	json.NewDecoder(w.Body).Decode(&userArrayLast)
	lastBytes, err := json.Marshal(userArrayLast)

	if !bytes.Equal(firstBytes, lastBytes) {
		t.Error("error, not equal")
	}
}

func TestInsertUsers(t *testing.T) {
	var err error
	var usersFromDb []User

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertUsers?first_name=TEST_NAME&second_name=TEST_SURNAME&date_of_birth=1900-01-21&role_id=777&password=password_test&login=TestUser", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM users WHERE first_name = 'TEST_NAME' AND second_name = 'TEST_SURNAME' AND login = 'TestUser' AND deleted IS NULL").Scan(&usersFromDb)

	if cap(usersFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(usersFromDb))
	}

	if usersFromDb[0].FirstName != "TEST_NAME" {
		t.Error("Wrong first_name")
	}

	if usersFromDb[0].SecondName != "TEST_SURNAME" {
		t.Error("Wrong second_name")
	}

	if usersFromDb[0].DateOfBirth.Format("2006-01-02") != "1900-01-21" {
		t.Error("Wrong date_of_birth")
	}

	if usersFromDb[0].RoleId != 777 {
		t.Error("Wrong role_id")
	}

	if usersFromDb[0].Login != "TestUser" {
		t.Error("Wrong login")
	}

	tmpPassHash := sha256.Sum256([]byte(r.FormValue("password")))
	tmpPassword := string(tmpPassHash[:])
	tmpPassword = fmt.Sprintf("%x", tmpPassword)
	if usersFromDb[0].Password != tmpPassword {
		t.Error("Wrong password")
	}

	if usersFromDb[0].Created.IsZero() {
		t.Error("Wrong create time")
	}

	testApi.Db.Exec("DELETE FROM users WHERE first_name = 'TEST_NAME' AND second_name = 'TEST_SURNAME' AND login = 'TestUser'")
}

func TestDeleteUsers(t *testing.T) {
	var err error
	var usersFromDb []User

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertUsers?first_name=TEST_NAME&second_name=TEST_SURNAME&date_of_birth=1900-01-21&role_id=777&password=password_test&login=TestUser", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM users WHERE role_id=777 AND deleted IS NULL").Scan(&usersFromDb)
	if cap(usersFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/deleteUsers?id="+strconv.Itoa(usersFromDb[0].Id)+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.DeleteUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM users WHERE role_id=777 AND deleted NOTNULL").Scan(&usersFromDb)
	if cap(usersFromDb) != 1 {
		t.Error("row didn't delete")
	}

	testApi.Db.Exec("DELETE FROM users WHERE role_id=777 AND deleted NOTNULL")
}

func TestUpdateUsers(t *testing.T) {
	var err error
	var usersFromDb []User

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertUsers?first_name=TEST_NAME&second_name=TEST_SURNAME&date_of_birth=1900-01-21&role_id=777&password=password_test&login=TestUser", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM users WHERE role_id=777 AND first_name = 'TEST_NAME' AND deleted IS NULL AND updated IS NULL").Scan(&usersFromDb)
	if cap(usersFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/updateUsers?id="+strconv.Itoa(usersFromDb[0].Id)+"&first_name=USER_N&second_name=USER_S&date_of_birth=1800-03-04&role_id=888&password=12344321qaz&login=TTT", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.UpdateUsers(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM users WHERE id=" + strconv.Itoa(usersFromDb[0].Id) + " AND updated NOTNULL").Scan(&usersFromDb)
	if cap(usersFromDb) != 1 {
		t.Error("row didn't update")
	}

	if usersFromDb[0].FirstName != "USER_N" {
		t.Error("Wrong updated first_name")
	}

	if usersFromDb[0].SecondName != "USER_S" {
		t.Error("Wrong updated second_name")
	}

	if usersFromDb[0].DateOfBirth.Format("2006-01-02") != "1800-03-04" {
		t.Error("Wrong updated date_of_birth")
	}

	if usersFromDb[0].RoleId != 888 {
		t.Error("Wrong updated role_id")
	}

	if usersFromDb[0].Login != "TTT" {
		t.Error("Wrong updated login")
	}

	tmpPassHash := sha256.Sum256([]byte("12344321qaz"))
	tmpPasswordStr := string(tmpPassHash[:])
	tmpPasword := fmt.Sprintf("%x", tmpPasswordStr)
	if usersFromDb[0].Password != tmpPasword {
		t.Error("Wrong updated pass")
	}

	if usersFromDb[0].Updated.IsZero() {
		t.Error("Wrong updated updated")
	}

	testApi.Db.Exec("DELETE FROM users WHERE id = " + strconv.Itoa(usersFromDb[0].Id) + " AND updated NOTNULL")
}
