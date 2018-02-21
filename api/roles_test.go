package api

import (
	"bytes"
	"encoding/json"
	"github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"
)

func TestGetRoles(t *testing.T) {
	var err error
	var roleArrayFirst, roleArrayLast []Role

	initApi(&testApi, t)
	defer testApi.Db.Close()

	testApi.Db.Raw("SELECT * FROM roles WHERE deleted IS NULL").Scan(&roleArrayFirst)

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/getRoles", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.GetRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
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

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertRoles?name=TEST&role_id=777&allowe_paths={/TEST1,/TEST2,/TEST3}", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
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

	if roleFromDb[0].Created.IsZero() {
		t.Error("Wrong creat time")
	}

	tmpAllowePathsStr := "{/TEST1,/TEST2,/TEST3}"
	tmpAllowePaths.Scan(tmpAllowePathsStr)

	if !reflect.DeepEqual(roleFromDb[0].AllowePaths, tmpAllowePaths) {
		t.Error("Wrong allowe_paths")
	}

	testApi.Db.Exec("DELETE FROM roles WHERE name = 'TEST' AND role_id = 777")
}

func TestDeleteRoles(t *testing.T) {
	var err error
	var rolesFromDb []Role

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertRoles?name=TEST&role_id=777&allowe_paths={/TEST1,/TEST2,/TEST3}", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM roles WHERE role_id=777 AND deleted IS NULL").Scan(&rolesFromDb)
	if cap(rolesFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/deleteRoles?id="+strconv.Itoa(rolesFromDb[0].Id)+"", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.DeleteRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM roles WHERE role_id=777 AND deleted NOTNULL").Scan(&rolesFromDb)
	if cap(rolesFromDb) != 1 {
		t.Error("row didn't delete")
	}

	testApi.Db.Exec("DELETE FROM roles WHERE role_id=777 AND deleted NOTNULL")
}

func TestUpdateRoles(t *testing.T) {
	var err error
	var rolesFromDb []Role

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertRoles?name=TEST&role_id=777&allow_paths={/TEST1,/TEST2,/TEST3}", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM roles WHERE role_id=777 AND deleted IS NULL").Scan(&rolesFromDb)
	if cap(rolesFromDb) != 1 {
		t.Error("row didn't insert")
	}

	w = httptest.NewRecorder()
	r, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/updateRoles?id="+strconv.Itoa(rolesFromDb[0].Id)+"&name=ROLEN&allowe_paths={/p1,/p2,/p3,/p4}&role_id=888", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.UpdateRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	testApi.Db.Raw("SELECT * FROM roles WHERE id=" + strconv.Itoa(rolesFromDb[0].Id) + " AND updated NOTNULL").Scan(&rolesFromDb)
	if cap(rolesFromDb) != 1 {
		t.Error("row didn't update")
	}

	if rolesFromDb[0].Name != "ROLEN" {
		t.Error("Wrong updated name")
	}

	var tmpAllowePaths pq.StringArray
	tmpAllowePathsStr := "{/p1,/p2,/p3,/p4}"
	tmpAllowePaths.Scan(tmpAllowePathsStr)

	if !reflect.DeepEqual(rolesFromDb[0].AllowePaths, tmpAllowePaths) {
		t.Error("Wrong allowe_paths")
	}

	if rolesFromDb[0].RoleId != 888 {
		t.Error("Wrong updated role_id")
	}

	if rolesFromDb[0].Updated.IsZero() {
		t.Error("wrong updated updated")
	}

	testApi.Db.Exec("DELETE FROM roles WHERE id = " + strconv.Itoa(rolesFromDb[0].Id) + " AND updated NOTNULL")
}

func TestGetRoleFromId(t *testing.T) {
	var err error
	var rolesFromDb Role

	initApi(&testApi, t)
	defer testApi.Db.Close()

	w := httptest.NewRecorder()
	r, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:8080/insertRoles?name=TEST&role_id=777&allow_paths={/TEST1,/TEST2,/TEST3}", nil)
	if err != nil {
		t.Fatal(err)
	}

	testApi.InsertRoles(w, r)
	if w.Code != http.StatusOK {
		t.Fatal(w.Code)
	}

	rolesFromDb, err = testApi.GetRoleFromRoleId(777)
	if rolesFromDb.Name != "TEST" {
		t.Error("Wrong Name, bad result")
	}

	testApi.Db.Exec("DELETE FROM roles WHERE role_id = 777")
}
