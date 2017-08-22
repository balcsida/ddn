package model

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/djavorszky/ddn/common/inet"
	"github.com/djavorszky/ddn/common/status"
	"github.com/djavorszky/notif"
	"github.com/djavorszky/sutils"
)

// DBRequest is used to represent JSON call about creating, dropping or importing databases
type DBRequest struct {
	ID           int    `json:"id"`
	DatabaseName string `json:"database_name"`
	DumpLocation string `json:"dumpfile_location"`
	Username     string `json:"username"`
	Password     string `json:"password"`
}

// ClientRequest is used to represent a JSON call between a client and the server
type ClientRequest struct {
	ConnectorIdentifier string `json:"connector_identifier"`
	RequesterEmail      string `json:"requester_email"`
	DBRequest
}

// RegisterRequest is used to represent a JSON call between the connector and the server.
// ID can be null if it's the initial registration, but must correspond to the connector's
// ID when unregistering
type RegisterRequest struct {
	ConnectorName string `json:"connector_name"`
	DBVendor      string `json:"dbvendor"`
	DBPort        string `json:"dbport"`
	DBAddr        string `json:"dbaddr"`
	DBSID         string `json:"dbsid"`
	ShortName     string `json:"short_name"`
	LongName      string `json:"long_name"`
	Version       string `json:"version"`
	Port          string `json:"port"`
	Addr          string `json:"address"`
}

// RegisterResponse is used as the response to the RegisterRequest
type RegisterResponse struct {
	ID      int    `json:"id"`
	Address string `json:"address"`
	Token   string `json:"token"`
}

// Connector is used to represent a DDN Connector.
type Connector struct {
	ID            int    `json:"id"`
	DBVendor      string `json:"vendor"`
	DBPort        string `json:"dbport"`
	DBAddr        string `json:"dbaddress"`
	DBSID         string `json:"sid"`
	ShortName     string `json:"connector"`
	LongName      string `json:"connector_long"`
	Identifier    string `json:"connector_identifier"`
	ConnectorPort string `json:"connector_port"`
	Version       string `json:"connector_version"`
	Address       string `json:"connector_address"`
	Token         string `json:"connector_token"`
	Up            bool   `json:"connector_up"`
}

// CreateDatabase sends a request to the connector to create a database.
func (c Connector) CreateDatabase(id int, dbname, dbuser, dbpass string) (string, error) {

	if ok := sutils.Present(dbname, dbuser, dbpass); !ok {
		return "", fmt.Errorf("asked to create database with missing values: dbname: %q, dbuser: %q, dbpass: %q", dbname, dbuser, dbpass)
	}

	dbreq := DBRequest{
		ID:           id,
		DatabaseName: dbname,
		Username:     dbuser,
		Password:     dbpass,
	}

	return c.executeAction(dbreq, "create-database")
}

// ImportDatabase starts the import on the connector.
func (c Connector) ImportDatabase(id int, dbname, dbuser, dbpass, dumploc string) (string, error) {
	if ok := sutils.Present(dbname, dbuser, dbpass, dumploc); !ok {
		return "", fmt.Errorf("asked to import database with missing values: dbname: %q, dbuser: %q, dbpass: %q, dumploc: %q", dbname, dbuser, dbpass, dumploc)
	}

	dbreq := DBRequest{
		ID:           id,
		DatabaseName: dbname,
		Username:     dbuser,
		Password:     dbpass,
		DumpLocation: dumploc,
	}

	return c.executeAction(dbreq, "import-database")
}

// DropDatabase sends a request to the connector to create a database.
func (c Connector) DropDatabase(id int, dbname, dbuser string) (string, error) {

	if ok := sutils.Present(dbname, dbuser); !ok {
		return "", fmt.Errorf("asked to create database with missing values: dbname: %q, dbuser: %q", dbname, dbuser)
	}

	dbreq := DBRequest{
		ID:           id,
		DatabaseName: dbname,
		Username:     dbuser,
	}

	return c.executeAction(dbreq, "drop-database")
}

func (c Connector) executeAction(dbreq DBRequest, endpoint string) (string, error) {
	dest := fmt.Sprintf("%s:%s/%s", c.Address, c.ConnectorPort, endpoint)

	if !strings.HasPrefix(dest, "http://") && !strings.HasPrefix(dest, "https://") {
		dest = fmt.Sprintf("http://%s", dest)
	}

	resp, err := notif.SndLoc(dbreq, dest)
	if err != nil && resp == "" {
		return "", fmt.Errorf("sending json message failed: %s", err.Error())
	}

	var respMsg inet.Message

	json.Unmarshal([]byte(resp), &respMsg)

	switch respMsg.Status {
	case status.Success, status.Accepted, status.Started, status.Created:
		return respMsg.Message, nil
	case status.MissingParameters:
		return "", fmt.Errorf("missing parameters from the request")
	case status.InvalidJSON:
		return "", fmt.Errorf("invalid JSON request")
	case status.CreateDatabaseFailed, status.ListDatabaseFailed, status.DropDatabaseFailed:
		return "", fmt.Errorf("connector issue: %s", respMsg.Message)
	default:
		return "", fmt.Errorf("executing action on endpoint %q failed: %s", endpoint, respMsg.Message)
	}
}
