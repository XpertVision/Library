package API

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
	"os/user"
)

func (a *API) LogOn(w http.ResponseWriter, r *http.Request) {
	var err error

	err = UniversalParseForm(&w, r)
	if err != nil {
		return
	}

	var user User
	var whereString string

	login := r.FormValue("login")
	WhereBlock("login", login, &whereString)

	pass := r.FormValue("password")
	hashPass := sha256.Sum256([]byte(pass))
	passStr := fmt.Sprintf("%x", string(hashPass[:]))

	WhereBlock("password", "'"+passStr+"'", &whereString)
	WhereBlock("deleted", "NULL", &whereString)

	query := "SELECT * FROM users WHERE " + whereString

	err = a.Db.Raw(query).Scan(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Authorization error"))
		return
	}

	if user.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong Login or Password"))
		return
	}

	var conn Connection
	genTime := time.Now()

	conn.UserId = user.Id
	conn.RoleId = user.RoleId
	conn.GenerateDate = genTime
	conn.Token = GetToken(user.Login, genTime)

	existConn, err := a.GetConnection(conn.Token)
	if err != nil && existConn.Id != 0 {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(existConn.Token))
		return
	}

	err = a.InsertConnection(conn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authirization error"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(conn.Token))
}

func (a *API) LogonHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		/*
				ПРИНИМАЮ TOKEN
			if token exist than get role_id
		*/

		/*if role_id cool for func -> run else error response
		 */
		var err error

		err = UniversalParseForm(&w, r)
		if err != nil {
			a.Log.Error("Parse form error")
			return
		}

		var conn Connection
		token := r.FormValue("token")

		err = a.Db.Raw("SELECT * FROM connections WHERE token = '" + token + "'").Scan(&conn).Error
		if err != nil {
			a.Log.Error("Inernal problems: func LogonHandler")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Inernal problems"))
		}

		if conn.GenerateDate.Minute()-time.Now().Minute() < -10 {
			var user User
			err = a.Db.Raw("SELECT * FROM users WHERE id = " + string(conn.UserId)).Scan(&user).Error
			if err != nil{
				a.Log.Error("Inernal problems: func LogonHandler | query from users table")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Inernal problems"))
			}

			genTime := time.Now()

			conn.UserId = user.Id
			conn.RoleId = user.RoleId
			conn.GenerateDate = genTime
			conn.Token = GetToken(user.Login, genTime)

			err = a.UpdateConnection(conn)
			if err != nil{
				a.Log.Error("Inernal problems: func LogonHandler | query update in connections")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("Inernal problems"))
			}
		}

		...
	}

	return http.HandlerFunc(fn)
}
