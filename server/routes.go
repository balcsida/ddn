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
		"register",
		http.MethodPost,
		"/register",
		register,
	},
	route{
		"unregister",
		http.MethodPost,
		"/unregister",
		unregister,
	},
	route{
		"heartbeat",
		http.MethodGet,
		"/heartbeat",
		heartbeat,
	},
	route{
		"alive",
		http.MethodGet,
		"/alive/{shortname:[a-zA-Z0-9-_]+}",
		alive,
	},
	route{
		"upd8",
		http.MethodPost,
		"/upd8",
		upd8,
	},
	route{
		"index",
		http.MethodGet,
		"/",
		index,
	},
	route{
		"create",
		http.MethodPost,
		"/create",
		createAction,
	},
	route{
		"createdb",
		http.MethodGet,
		"/createdb",
		createdb,
	},
	route{
		"import",
		http.MethodPost,
		"/import",
		importAction,
	},
	route{
		"prepimport",
		http.MethodPost,
		"/prepimport",
		prepImportAction,
	},
	route{
		"importdb",
		http.MethodGet,
		"/importdb",
		importdb,
	},
	route{
		"fileimport",
		http.MethodGet,
		"/fileimport",
		fileimport,
	},
	route{
		"srvimport",
		http.MethodGet,
		"/srvimport",
		srvimport,
	},
	route{
		"browse",
		http.MethodGet,
		"/browse/{loc:[0-9a-zA-Z-_./ ]+}",
		browse,
	},
	route{
		"browse",
		http.MethodGet,
		"/browse",
		browseroot,
	},
	route{
		"login",
		http.MethodPost,
		"/login",
		login,
	},
	route{
		"logout",
		http.MethodGet,
		"/logout",
		logout,
	},
	route{
		"extend",
		http.MethodGet,
		"/extend/{id:[0-9]+}",
		extend,
	},
	route{
		"drop",
		http.MethodGet,
		"/drop/{id:[0-9]+}",
		drop,
	},
	route{
		"portalext",
		http.MethodGet,
		"/portalext/{id:[0-9]+}",
		portalext,
	},
	route{
		"recreate",
		http.MethodGet,
		"/recreate/{id:[0-9]+}",
		recreate,
	},
	route{
		"api",
		http.MethodGet,
		"/api",
		apiPage,
	},
	route{
		"api/create",
		http.MethodPost,
		"/api/create",
		apiCreate,
	},
	route{
		"api/list",
		http.MethodGet,
		"/api/list",
		apiList,
	},
	route{
		"api/listAgents",
		http.MethodGet,
		"/api/list-agents",
		apiListAgents,
	},

	route{
		"api/listDatabases",
		http.MethodPost,
		"/api/list-databases",
		apiListDatabases,
	},
	/*
		route{
			"api/agents/byName",
			http.MethodGet,
			"/api/agents/{shortname:[a-zA-Z0-9-_]+}",
			apiAgentByName,
		},
	*/
	route{
		"api/visibility/",
		http.MethodGet,
		"/api/visibility/{id:[0-9]+}/{visibility:public|private}",
		apiSetVisibility,
	},
	route{
		"safe2restart",
		http.MethodGet,
		"/api/safe2restart",
		apiSafe2Restart,
	},
	route{
		"api/save-subscription",
		http.MethodPost,
		"/api/save-subscription",
		apiSaveSubscription,
	},
	route{
		"api/remove-subscription",
		http.MethodPost,
		"/api/remove-subscription",
		apiRemoveSubscription,
	},
	route{
		"/api/dbaccess/",
		http.MethodGet,
		"/api/dbaccess/{requester:[a-zA-Z0-9-_.@]+}/{agent:[a-zA-Z0-9-_]+}/{dbname:[a-zA-Z0-9-_]+}",
		apiDBAccess,
	},
	// v2 APIs
	route{
		"/api/agents",
		http.MethodGet,
		"/api/agents",
		getAPIAgents,
	},
	route{
		"/api/agents/active",
		http.MethodGet,
		"/api/agents/active",
		getAPIActiveAgents,
	},
	route{
		"api/agents/$agent-name",
		http.MethodGet,
		"/api/agents/{agent:[a-zA-Z0-9-_]+}",
		getAPIAgentByName,
	},
	route{
		"api/databases",
		http.MethodGet,
		"/api/databases",
		getAPIDatabases,
	},
	route{
		"api/databases/id",
		http.MethodGet,
		"/api/databases/{id:[0-9]+}",
		getAPIDatabaseByID,
	},
	route{
		"api/databases/agent/dbname",
		http.MethodGet,
		"/api/databases/{agent:[a-zA-Z][a-zA-Z0-9-_]+}/{dbname:[a-zA-Z0-9_]+}",
		getAPIDatabaseByAgentDBName,
	},
	route{
		"api/databases/id",
		http.MethodDelete,
		"/api/databases/{id:[0-9]+}",
		dropAPIDatabaseByID,
	},
	route{
		"api/databases/create",
		http.MethodPost,
		"/api/databases/create",
		createAPIDB,
	},
	route{
		"api/databases/import",
		http.MethodPost,
		"/api/databases/import",
		importAPIDB,
	},
	route{
		"api/databases/id/recreate",
		http.MethodPut,
		"/api/databases/{id:[0-9]+}/recreate",
		recreateAPIDB,
	},
	route{
		"api/browse",
		http.MethodGet,
		"/api/browse",
		browseAPI,
	},
	route{
		"api/browse/loc",
		http.MethodGet,
		"/api/browse/{loc:[0-9a-zA-Z-_./ ]+}",
		browseAPI,
	},
	route{
		"api/databases/visibility",
		http.MethodPut,
		"/api/databases/{id:[0-9]+}/visibility/{visibility:public|private}",
		apiSetVisibility,
	},
	route{
		"api/databases/expiry",
		http.MethodPut,
		"/api/databases/{id:[0-9]+}/expiry/extend/{amount:[0-9]+}/{unit:days|months|years}",
		apiExtendExpiry,
	},
	route{
		"/api/databases/id/accessinfo/",
		http.MethodGet,
		"/api/databases/{id:[0-9]+}/accessinfo",
		apiAccessInfoByID,
	},
	route{
		"/api/dbaccess/",
		http.MethodGet,
		"/api/databases/{agent:[a-zA-Z][a-zA-Z0-9-_]+}/{dbname:[a-zA-Z0-9-_]+}/accessinfo",
		apiAccessInfoByAgentDB,
	},
	route{
		"api/loglevel",
		http.MethodPut,
		"/api/loglevel/{level:[a-zA-Z]+}",
		apiSetLogLevel,
	},
}
