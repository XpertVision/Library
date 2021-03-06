package wrappers

import (
	"net/http"

	"github.com/XpertVision/Library/api"
	"github.com/justinas/alice"
)

//HandleAll func is wrapper where makes all handle connection and middleware
func HandleAll(api *api.API) {

	wrapperLogOn := alice.New(api.ParseHandler)
	http.Handle("/logon", wrapperLogOn.ThenFunc(api.LogOn))

	wrapper := alice.New(api.ParseHandler, api.LogonHandler)

	http.Handle("/getBooks", wrapper.ThenFunc(api.GetBooks))
	http.Handle("/getRoles", wrapper.ThenFunc(api.GetRoles))
	http.Handle("/getUsers", wrapper.ThenFunc(api.GetUsers))

	http.Handle("/updateUsers", wrapper.ThenFunc(api.UpdateUsers))
	http.Handle("/updateRoles", wrapper.ThenFunc(api.UpdateRoles))
	http.Handle("/updateBooks", wrapper.ThenFunc(api.UpdateBooks))

	http.Handle("/insertBooks", wrapper.ThenFunc(api.InsertBooks))
	http.Handle("/insertUsers", wrapper.ThenFunc(api.InsertUsers))
	http.Handle("/insertRoles", wrapper.ThenFunc(api.InsertRoles))

	http.Handle("/deleteRoles", wrapper.ThenFunc(api.DeleteRoles))
	http.Handle("/deleteBooks", wrapper.ThenFunc(api.DeleteBooks))
	http.Handle("/deleteUsers", wrapper.ThenFunc(api.DeleteUsers))
}
