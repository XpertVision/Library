package api

import (
	//"strconv"
	"strconv"
	"testing"
	"time"
)

func TestGetConnectionFromId(t *testing.T) {
	var err error
	var connectionToDb, connectionFromFunc Connection
	var connectionFromDb []Connection

	initApi(&testApi, t)
	defer testApi.Db.Close()

	connectionToDb.UserId = 777
	connectionToDb.RoleId = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testApi.InsertConnection(connectionToDb)

	testApi.Db.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	connectionFromFunc, err = testApi.GetConnectionFromId(connectionFromDb[0].UserId)
	if err != nil {
		t.Error("Error: ", err)
	}

	if connectionFromFunc.Id != connectionFromDb[0].Id {
		t.Error("problem with id: ", connectionFromFunc.Id)
	}

	if connectionFromFunc.UserId != connectionFromDb[0].UserId {
		t.Error("problem with user_id: ", connectionFromFunc.UserId)
	}

	if connectionFromFunc.RoleId != connectionFromDb[0].RoleId {
		t.Error("problem with role_id: ", connectionFromFunc.RoleId)
	}

	if connectionFromFunc.Token != connectionFromDb[0].Token {
		t.Error("problem with token: ", connectionFromFunc.Token)
	}

	if connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05") != connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("problem wuth generate_date: ", connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05"))
	}

	testApi.Db.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionFromFunc.Id))
}

func TestGetConnectionFromToken(t *testing.T) {
	var err error
	var connectionToDb, connectionFromFunc Connection
	var connectionFromDb []Connection

	initApi(&testApi, t)
	defer testApi.Db.Close()

	connectionToDb.UserId = 777
	connectionToDb.RoleId = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testApi.InsertConnection(connectionToDb)

	testApi.Db.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	connectionFromFunc, err = testApi.GetConnectionFromToken(connectionFromDb[0].Token)
	if err != nil {
		t.Error("Error: ", err)
	}

	if connectionFromFunc.Id != connectionFromDb[0].Id {
		t.Error("problem with id")
	}

	if connectionFromFunc.UserId != connectionFromDb[0].UserId {
		t.Error("problem with user_id")
	}

	if connectionFromFunc.RoleId != connectionFromDb[0].RoleId {
		t.Error("problem with role_id")
	}

	if connectionFromFunc.Token != connectionFromDb[0].Token {
		t.Error("problem with token")
	}

	if connectionFromFunc.GenerateDate.Format("2006-01-02 15:04:05") != connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("problem wuth generate_date")
	}

	testApi.Db.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionFromFunc.Id))
}

func TestInsertConnection(t *testing.T) {
	var err error
	var connectionToDb Connection
	var connectionFromDb []Connection

	initApi(&testApi, t)
	defer testApi.Db.Close()

	connectionToDb.UserId = 777
	connectionToDb.RoleId = 777
	connectionToDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionToDb.Token = "test_token"

	testApi.InsertConnection(connectionToDb)

	testApi.Db.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionFromDb)

	if cap(connectionFromDb) != 1 {
		t.Error("too much rows found, capacity: ", cap(connectionFromDb))
	}

	if connectionFromDb[0].UserId != connectionToDb.UserId {
		t.Error("Wrong user_id")
	}

	if connectionFromDb[0].RoleId != connectionToDb.RoleId {
		t.Error("Wrong role_id")
	}

	if connectionFromDb[0].GenerateDate.Format("2006-01-02 15:04:05") != connectionToDb.GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("Wrong generate_date:::FromDB: ", connectionFromDb[0].GenerateDate, " :::ToDb: ", connectionToDb.GenerateDate)
	}

	if connectionFromDb[0].Token != connectionToDb.Token {
		t.Error("Wrong token")
	}

	testApi.Db.Exec("DELETE FROM connections WHERE role_id = 777 AND token = 'test_token'")
}

func TestUpdateConnection(t *testing.T) {
	var err error
	var connectionInDb, connectingForUp Connection

	initApi(&testApi, t)
	defer testApi.Db.Close()

	connectionInDb.UserId = 777
	connectionInDb.RoleId = 777
	connectionInDb.GenerateDate, err = time.Parse("2006-01-02 15:04:05", "1900-01-02 15:15:15")
	if err != nil {
		t.Fatal("parse date error: ", err)
	}
	connectionInDb.Token = "test_token"

	testApi.InsertConnection(connectionInDb)

	testApi.Db.Raw("SELECT * FROM connections WHERE user_id = 777").Scan(&connectionInDb)

	connectingForUp.Id = connectionInDb.Id
	connectingForUp.UserId = 888
	connectingForUp.RoleId = 888
	connectingForUp.Token = "update_token"
	connectingForUp.GenerateDate = time.Now()

	testApi.UpdateConnection(connectingForUp)

	testApi.Db.Raw("SELECT * FROM connections WHERE id = " + strconv.Itoa(connectingForUp.Id)).Scan(&connectionInDb)

	if connectingForUp.Id != connectionInDb.Id {
		t.Error("Wrong id")
	}

	if connectingForUp.UserId != connectionInDb.UserId {
		t.Error("Wrong user_id")
	}

	if connectingForUp.RoleId != connectionInDb.RoleId {
		t.Error("Wrong role_id")
	}

	if connectingForUp.Token != connectionInDb.Token {
		t.Error("Wrong token")
	}

	if connectingForUp.GenerateDate.Format("2006-01-02 15:04:05") != connectionInDb.GenerateDate.Format("2006-01-02 15:04:05") {
		t.Error("Wrong generate_date")
	}

	testApi.Db.Exec("DELETE FROM connections WHERE id = " + strconv.Itoa(connectionInDb.Id))
}

func TestDeleteConnection(t *testing.T) {
	var conn Connection

	initApi(&testApi, t)
	defer testApi.Db.Close()

	conn.UserId = 777
	conn.RoleId = 777
	conn.GenerateDate = time.Now()
	conn.Token = "test_token"

	testApi.InsertConnection(conn)

	testApi.Db.Raw("SELECT * FROM connections WHERE token = '" + conn.Token + "'").Scan(&conn)

	if conn.Id == 0 {
		t.Fatal("Insert was incorrect")
	}

	testApi.DeleteConnection(conn.Token)

	conn.Id = 0
	testApi.Db.Raw("SELECT * FROM connections WHERE token = '" + conn.Token + "'").Scan(&conn)

	if conn.Id != 0 {
		t.Fatal("Delete was incorrect", conn.Id)
	}
}
