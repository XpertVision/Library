package api

import (
	//"strconv"
	"strconv"
	"testing"
	"time"
)

func TestGetConnectionFromID(t *testing.T) {
	var err error
	var connectionToDb, connectionFromFunc Connection
	var connectionFromDb []Connection

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	connectionToDb.UserID = 777
	connectionToDb.RoleID = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testAPI.InsertConnection(connectionToDb)

	testAPI.DB.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	connectionFromFunc, err = testAPI.GetConnectionFromID(connectionFromDb[0].UserID)
	if err != nil {
		t.Error("Error: ", err)
	}

	if connectionFromFunc.ID != connectionFromDb[0].ID {
		t.Error("problem with id: ", connectionFromFunc.ID)
	}

	if connectionFromFunc.UserID != connectionFromDb[0].UserID {
		t.Error("problem with user_id: ", connectionFromFunc.UserID)
	}

	if connectionFromFunc.RoleID != connectionFromDb[0].RoleID {
		t.Error("problem with role_id: ", connectionFromFunc.RoleID)
	}

	if connectionFromFunc.Token != connectionFromDb[0].Token {
		t.Error("problem with token: ", connectionFromFunc.Token)
	}

	if connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05") != connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("problem wuth generate_date: ", connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05"))
	}

	testAPI.DB.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionFromFunc.ID))
}

func TestGetConnectionFromToken(t *testing.T) {
	var err error
	var connectionToDb, connectionFromFunc Connection
	var connectionFromDb []Connection

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	connectionToDb.UserID = 777
	connectionToDb.RoleID = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testAPI.InsertConnection(connectionToDb)

	testAPI.DB.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	connectionFromFunc, err = testAPI.GetConnectionFromToken(connectionFromDb[0].Token)
	if err != nil {
		t.Error("Error: ", err)
	}

	if connectionFromFunc.ID != connectionFromDb[0].ID {
		t.Error("problem with id")
	}

	if connectionFromFunc.UserID != connectionFromDb[0].UserID {
		t.Error("problem with user_id")
	}

	if connectionFromFunc.RoleID != connectionFromDb[0].RoleID {
		t.Error("problem with role_id")
	}

	if connectionFromFunc.Token != connectionFromDb[0].Token {
		t.Error("problem with token")
	}

	if connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05") != connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("problem wuth generate_date")
	}

	testAPI.DB.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionFromFunc.ID))
}

func TestInsertConnection(t *testing.T) {
	var err error
	var connectionToDb Connection
	var connectionFromDb []Connection

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	connectionToDb.UserID = 777
	connectionToDb.RoleID = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testAPI.InsertConnection(connectionToDb)

	testAPI.DB.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	if connectionFromDb[0].UserID != connectionToDb.UserID {
		t.Error("Wrong user_id")
	}

	if connectionFromDb[0].RoleID != connectionToDb.RoleID {
		t.Error("Wrong role_id")
	}

	if connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") != connectionToDb.GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("Wrong generate_date:::FromDB: ", connectionFromDb[0].GenerateDate, " :::ToDb: ", connectionToDb.GenerateDate)
	}

	if connectionFromDb[0].Token != connectionToDb.Token {
		t.Error("Wrong token")
	}

	testAPI.DB.Exec("DELETE FROM connections WHERE role_id = 777 AND token = 'test_token'")
}

func TestUpdateConnection(t *testing.T) {
	var err error
	var connectionInDb, connectingForUp Connection

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	connectionInDb.UserID = 777
	connectionInDb.RoleID = 777
	connectionInDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionInDb.Token = "test_token"

	testAPI.InsertConnection(connectionInDb)

	testAPI.DB.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionInDb)

	connectingForUp.ID = connectionInDb.ID
	connectingForUp.UserID = 888
	connectingForUp.RoleID = 888
	connectingForUp.Token = "update_token"
	connectingForUp.GenerateDate = time.Now()

	testAPI.UpdateConnection(connectingForUp)

	testAPI.DB.Raw("SELECT * FROM connections WHERE id = " + strconv.Itoa(connectingForUp.ID)).Scan(&connectionInDb)

	if connectingForUp.ID != connectionInDb.ID {
		t.Error("Wrong id")
	}

	if connectingForUp.UserID != connectionInDb.UserID {
		t.Error("Wrong user_id")
	}

	if connectingForUp.RoleID != connectionInDb.RoleID {
		t.Error("Wrong role_id")
	}

	if connectingForUp.Token != connectionInDb.Token {
		t.Error("Wrong token")
	}

	if connectingForUp.GenerateDate.Format("2006-01-02 15:04:05") != connectionInDb.GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("Wrong generate_date")
	}

	testAPI.DB.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionInDb.ID))
}

func TestDeleteConnection(t *testing.T) {
	var err error
	var conn Connection

	initAPI(&testAPI, t)
	defer testAPI.DB.Close()

	conn.UserID = 777
	conn.RoleID = 777
	conn.GenerateDate = time.Now()
	conn.Token = "test_token"

	err = testAPI.InsertConnection(conn)
	if err != nil {
		t.Fatal(err)
	}

	testAPI.DB.Raw("SELECT * FROM connections WHERE token = '" + conn.Token + "'").Scan(&conn)

	if conn.ID == 0 {
		t.Fatal("Insert was incorrect")
	}

	testAPI.DeleteConnection(conn.Token)

	conn.ID = 0
	testAPI.DB.Raw("SELECT * FROM connections WHERE token = '" + conn.Token + "'").Scan(&conn)

	if conn.ID != 0 {
		t.Fatal("Delete was incorrect", conn.ID)
	}
}
