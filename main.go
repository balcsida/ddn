package main

import (
	"fmt"
	"log"

	"net/http"

	"os/user"

	"os"

	"strings"

	"github.com/BurntSushi/toml"
)

var (
	properties string
	conf       Config
	db         Database
	port       string
	usr        *user.User
)

func main() {
	var err error

	usr, err = user.Current()
	if err != nil {
		log.Fatal("Couldn't get default user.")
	}

	properties, err := checkProps()
	if err != nil {
		log.Println("Couldn't find properties file, generating one")
		file, err := generateProps()
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("Generated '%s' with dummy values next to executable. Please update it with real values and restart the connector", file)
	}

	if _, err := toml.DecodeFile(properties, &conf); err != nil {
		log.Fatal(err)
	}

	log.Println("Starting with properties:")
	log.Println("Vendor:\t\t", conf.Vendor)
	log.Println("Version:\t\t", conf.Version)
	log.Println("Database port:\t", conf.DBPort)
	log.Println("Connector port:\t", conf.ConnectorPort)
	log.Println("Executable:\t\t", conf.Exec)
	log.Println("Username:\t\t", conf.User)
	log.Println("Password:\t\t ******")
	log.Println("Master address:\t", conf.MasterAddress)

	if _, err = os.Stat(conf.Exec); os.IsNotExist(err) {
		log.Fatalf("Database executable '%s' doesn't exist.", conf.Exec)
	}

	switch strings.ToLower(conf.Vendor) {
	case "mysql":
		db = new(mysql)
	default:
		log.Fatal("Database vendor not recognized.")
	}

	err = db.Connect(conf.Vendor, conf.User, conf.Password, conf.DBPort)
	if err != nil {
		log.Fatal("Could not establish database connection:\n\t\t", err.Error())
	}
	defer db.Close()

	log.Println("Database connection established")

	log.Println("Starting to listen on port", conf.ConnectorPort)

	port = fmt.Sprintf(":%s", conf.ConnectorPort)

	log.Fatal(http.ListenAndServe(port, Router()))

}
