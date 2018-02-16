package wrappers

import (
	"../api"
	"github.com/justinas/alice"
	"net/http"
)

func HandleAll(api *api.API) {
	http.HandleFunc("/logon", api.LogOn)

	wrapper := alice.New(api.LogonHandler)
	http.Handle("/getBooks", wrapper.ThenFunc(api.GetBooks))
	//http.HandleFunc("/getBooks", api.GetBooks)
	http.HandleFunc("/getRoles", api.GetRoles)
	http.HandleFunc("/getUsers", api.GetUsers)

	http.HandleFunc("/updateUsers", api.UpdateUsers)
	http.HandleFunc("/updateRoles", api.UpdateRoles)
	http.HandleFunc("/updateBooks", api.UpdateBooks)

	http.HandleFunc("/insertBooks", api.InsertBooks)
	http.HandleFunc("/insertUsers", api.InsertUsers)
	http.HandleFunc("/insertRoles", api.InsertRoles)

	http.HandleFunc("/deleteRoles", api.DeleteRoles)
	http.HandleFunc("/deleteBooks", api.DeletetBooks)
	http.HandleFunc("/deleteUsers", api.DeletetUsers)

}
