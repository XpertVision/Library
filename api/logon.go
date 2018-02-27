package api

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

//LogOn func checks exist user/pass pare or not, if exist - generates and return token, if logined yet - return current token
func (a *API) LogOn(w http.ResponseWriter, r *http.Request) {
	var err error

	var user, userTmp User

	userTmp.Login = r.FormValue("login")

	pass := r.FormValue("password")
	hashPass := sha256.Sum256([]byte(pass))
	userTmp.Password = fmt.Sprintf("%x", string(hashPass[:]))

	//Search user in db
	err = a.DB.Where(&userTmp).Find(&user).Error
	if err != nil {
		a.Log.Error("problem with select query | Error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("ERROR"))
		return
	}

	//If wrong user
	if user.ID == 0 {
		a.Log.Info("Incorrect user")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Wrong Login or Password"))
		return
	}

	//Return token if exist for this user
	conn, err := a.GetConnectionFromID(user.ID)
	if err == nil && conn.ID != 0 {
		a.Log.Info("authorizated yet")
		w.WriteHeader(http.StatusAccepted)
		w.Write([]byte(conn.Token))
		return
	}

	genTime := time.Now()

	conn.UserID = user.ID
	conn.RoleID = user.RoleID
	conn.GenerateDate = genTime
	conn.Token = getToken(user.Login, genTime)

	//Loginig user
	err = a.InsertConnection(conn)
	if err != nil {
		a.Log.Error("problem with insert new connection to db | Error: ", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Authirization error"))
		return
	}

	//return token
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte(conn.Token))
}

//ParseHandler func is middleware for checking form valids
func (a *API) ParseHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			a.Log.Error("Invalid form | Error: ", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("ERROR"))
			return
		}

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

//LogonHandler func is middleware for cheching valids of token and access role for users
func (a *API) LogonHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var err error

		var conn Connection
		token := r.FormValue("token")

		//If invalid token
		err = a.DB.Find(&conn).Where("token = ?", token).Error
		if err != nil {
			a.Log.Error("problem with select query | Error: ", err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("ERROR"))
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
			a.Log.Error("update token time error | Error: ", err)
		}

		//Finding role_id
		role, err := a.GetRoleFromRoleID(conn.RoleID)
		if err != nil {
			a.Log.Error("get role from role id error | Errore: ", err)
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
			a.Log.Info("Inccorect role for func")
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Access Denied"))
			return
		}

		//Call target-function
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
