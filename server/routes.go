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
		"echo",
		"POST",
		"/echo",
		echo,
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
}
