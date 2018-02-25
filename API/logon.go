package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

func (a *API) LogOn(w http.ResponseWriter, r *http.Request) {
	var err error

	var user, userTmp User

	userTmp.Login = r.FormValue("login")

	pass := r.FormValue("password")
	hashPass := sha256.Sum256([]byte(pass))
	userTmp.Password = fmt.Sprintf("%x", string(hashPass[:]))

	//Search user in db
	err = a.Db.Where(&userTmp).Find(&user).Error
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
	conn.Token = getToken(user.Login, genTime)

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

func (a *API) ParseHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("BAD REQUEST: Parse form error"))
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func (a *API) LogonHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error

		var conn Connection
		token := r.FormValue("token")

		//If invalid token
		err = a.Db.Find(&conn).Where("token = ?", token).Error
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
