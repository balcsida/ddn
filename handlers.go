package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/djavorszky/ddn/inet"
	"github.com/djavorszky/ddn/model"
	"github.com/djavorszky/notif"
	"github.com/djavorszky/sutils"
)

// index should display whenever someone visits the main page.

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to the index!")
}

func createDatabase(w http.ResponseWriter, r *http.Request) {
	var (
		dbreq model.DBRequest
		msg   inet.Message
	)

	err := json.NewDecoder(r.Body).Decode(&dbreq)
	if err != nil {
		log.Printf("couldn't decode json request: %s", err.Error())

		inet.SendResponse(w, inet.ErrorJSONResponse(err))
		return
	}

	if ok := sutils.Present(db.RequiredFields(dbreq, createDB)...); !ok {
		inet.SendResponse(w, inet.InvalidResponse())
		return
	}

	err = db.CreateDatabase(dbreq)
	if err != nil {
		log.Printf("creating database failed: %s", err.Error())
		msg.Status = http.StatusInternalServerError
		msg.Message = err.Error()
	} else {
		msg.Status = http.StatusOK
		msg.Message = "Successfully created the database and user!"
	}

	inet.SendResponse(w, msg)
}

// listDatabase lists the supervised databases in a JSON format
func listDatabases(w http.ResponseWriter, r *http.Request) {
	var (
		msg inet.ListMessage
		err error
	)

	msg.Status = http.StatusOK
	msg.Message, err = db.ListDatabase()
	if err != nil {
		log.Printf("listing databases failed: %s", err.Error())

		inet.SendResponse(w, inet.ErrorResponse())
		return
	}

	inet.SendResponse(w, msg)
}

// echo echoes whatever it receives (as JSON) to the log.
func echo(w http.ResponseWriter, r *http.Request) {
	var msg notif.Msg

	err := json.NewDecoder(r.Body).Decode(&msg)
	if err != nil {
		log.Printf("couldn't decode json request: %s", err.Error())

		inet.SendResponse(w, inet.ErrorJSONResponse(err))
		return
	}

	log.Printf("%+v", msg)
}

// dropDatabase will drop the named database with its tablespace and user
func dropDatabase(w http.ResponseWriter, r *http.Request) {
	var (
		dbreq model.DBRequest
		msg   inet.Message
	)

	err := json.NewDecoder(r.Body).Decode(&dbreq)
	if err != nil {
		log.Printf("couldn't drop database: %s", err.Error())

		inet.SendResponse(w, inet.ErrorJSONResponse(err))
		return
	}

	if ok := sutils.Present(db.RequiredFields(dbreq, dropDB)...); !ok {
		inet.SendResponse(w, inet.InvalidResponse())
		return
	}

	err = db.DropDatabase(dbreq)
	if err != nil {
		log.Printf("dropping database failed: %s", err.Error())

		msg.Status = http.StatusInternalServerError
		msg.Message = err.Error()
	} else {
		msg.Status = http.StatusOK
		msg.Message = "Successfully dropped the database and user!"
	}

	inet.SendResponse(w, msg)

}

// importDatabase will import the specified dumpfile to the database
// creating the database, tablespace and user
func importDatabase(w http.ResponseWriter, r *http.Request) {
	var (
		dbreq model.DBRequest
		msg   inet.Message
	)

	err := json.NewDecoder(r.Body).Decode(&dbreq)
	if err != nil {
		log.Printf("couldn't decode json request: %s", err.Error())

		inet.SendResponse(w, inet.ErrorJSONResponse(err))
		return
	}

	if ok := sutils.Present(db.RequiredFields(dbreq, importDB)...); !ok {
		inet.SendResponse(w, inet.InvalidResponse())
		return
	}

	if exists := inet.FileExists(dbreq.DumpLocation); exists == false {
		msg.Status = http.StatusNotFound
		msg.Message = "Specified file doesn't exist or is not reachable."

		inet.SendResponse(w, msg)
		return
	}

	err = db.CreateDatabase(dbreq)
	if err != nil {
		log.Printf("creating database failed: %s", err.Error())

		msg.Status = http.StatusInternalServerError
		msg.Message = err.Error()

		inet.SendResponse(w, msg)
		return
	}

	msg.Status = http.StatusOK
	msg.Message = "Understood request, starting import process."

	inet.SendResponse(w, msg)

	go startImport(dbreq)
}

func whoami(w http.ResponseWriter, r *http.Request) {
	info := make(map[string]string)

	info["database-vendor"] = conf.Vendor
	info["database-version"] = conf.Version
	info["connector-version"] = version

	duration := time.Since(startup)

	// Round to milliseconds.
	info["connector-uptime"] = fmt.Sprintf("%s", duration-(duration%time.Millisecond))

	var msg inet.MapMessage

	msg.Status = http.StatusOK
	msg.Message = info

	inet.SendResponse(w, msg)
}

func heartbeat(w http.ResponseWriter, r *http.Request) {
	var msg inet.Message

	msg.Status = http.StatusOK
	msg.Message = "Still alive"

	err := db.Alive()
	if err != nil {
		log.Printf("database dead: %s", err.Error())
		msg = inet.ErrorResponse()
	}

	inet.SendResponse(w, msg)
}
