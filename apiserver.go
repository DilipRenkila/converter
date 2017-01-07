/*
apiserver is responsible for serving the currency conversion API on port 8080.
*/
package main

import (
	"log"
	"net/http"
	"github.com/diliprenkila/converter/converter"
	"github.com/gorilla/mux"
)

/*
Route runs a mux Router.
*/
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

/*
Routes an array of type Route.
*/
type Routes []Route

/*
Define the "GET /convert" route.
*/
var routes = Routes{
	Route{
		"convert",
		"GET",
		"/convert",
		converter.ConvertCurrency,
	},
}

/*
NewRouter creates a mux router for the routes define above.
*/
func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	return router
}

func main() {

	router := NewRouter()
	// serves http on port 8080.
	log.Fatal(http.ListenAndServe(":8080", router))
}
