package main

import "net/http"

type route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes contains all available routes
type Routes []route

var routes = Routes{
	route{
		"listConnectors",
		"GET",
		"/list-connectors",
		listConnectors,
	},
	route{
		"register",
		"POST",
		"/register",
		register,
	},
	route{
		"unregister",
		"POST",
		"/unregister",
		unregister,
	},
	route{
		"alive",
		"GET",
		"/alive",
		alive,
	},
	route{
		"upd8",
		"POST",
		"/upd8",
		upd8,
	},
	route{
		"index",
		"GET",
		"/",
		index,
	},
	route{
		"create",
		"POST",
		"/create",
		createAction,
	},
	route{
		"createdb",
		"GET",
		"/createdb",
		createdb,
	},
	route{
		"import",
		"POST",
		"/import",
		importAction,
	},
	route{
		"importdb",
		"GET",
		"/importdb",
		importdb,
	},
	route{
		"login",
		"POST",
		"/login",
		login,
	},
	route{
		"logout",
		"GET",
		"/logout",
		logout,
	},
	route{
		"extend",
		"GET",
		"/extend/{id:[0-9]+}",
		extend,
	},
	route{
		"drop",
		"GET",
		"/drop/{id:[0-9]+}",
		drop,
	},
	route{
		"portalext",
		"GET",
		"/portalext/{id:[0-9]+}",
		portalext,
	},
	route{
		"api/create",
		"POST",
		"/api/create",
		apiCreate,
	},
	route{
		"api/list",
		"POST",
		"/api/list",
		apiList,
	},
	route{
		"api/import",
		"POST",
		"/api/import",
		apiImport,
	},
}
