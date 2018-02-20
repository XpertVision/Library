package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
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

	//Search user in db
	err = a.Db.Raw(query).Scan(&user).Error
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Authorization error"))
		return
	}

	//If wrong user
	if user.Id == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong Login or Password"))
		return
	}

	//Return token if exist for this user
	existConn, err := a.GetConnectionFromId(user.Id)
	if err == nil && existConn.Id != 0 {
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(existConn.Token))
		return
	}

	var conn Connection
	genTime := time.Now()

	conn.UserId = user.Id
	conn.RoleId = user.RoleId
	conn.GenerateDate = genTime
	conn.Token = GetToken(user.Login, genTime)

	//Loginig user
	err = a.InsertConnection(conn)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authirization error"))
		return
	}

	//return token
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(conn.Token))
}

func (a *API) LogonHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error

		err = UniversalParseForm(&w, r)
		if err != nil {
			a.Log.Error("Parse form error")
			return
		}

		var conn Connection
		token := r.FormValue("token")

		//If invalid token
		err = a.Db.Raw("SELECT * FROM connections WHERE token = '" + token + "'").Scan(&conn).Error
		if err != nil {
			a.Log.Error("Inernal problems: func LogonHandler")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Incorrect user"))
			return
		}

		//If token too old
		if (conn.GenerateDate.Unix()-(time.Now().Unix()+7200))/60 < -10 {
			err = a.DeleteConnection(conn.Token)
			if err != nil {
				a.Log.Error(err.Error())
			}

			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Session TimeOut"))
			return
		}

		//If all ok with token, update his generate date
		conn.GenerateDate = time.Now()
		err = a.UpdateConnection(conn)
		if err != nil {
			a.Log.Error("update token time error")
		}

		//Finding role_id
		role, err := a.GetRoleFromRoleId(conn.RoleId)
		if err != nil {
			a.Log.Error(err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Access Denied"))
			return
		}

		//Knowing access rule for user anf compare with target-function
		var accessOk bool
		for _, strRole := range role.AllowePaths {
			if strRole == r.URL.Path {
				accessOk = true
				break
			}
		}

		//If access role incorrect for this target-func
		if !accessOk {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Access Denied"))
			return
		}

		//Call target-function
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
