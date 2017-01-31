package main

import (
	"log"
	"runtime"

	"github.com/djavorszky/prompter"
)

// Config to hold the database server information
type Config struct {
	Vendor            string `toml:"vendor"`
	Version           string `toml:"version"`
	Exec              string `toml:"executable"`
	DBPort            string `toml:"dbport"`
	DBAddress         string `toml:"dbaddress"`
	SID               string `toml:"oracle-sid"`
	ConnectorPort     string `toml:"connectorPort"`
	User              string `toml:"username"`
	Password          string `toml:"password"`
	DefaultTablespace string `toml:"default-tablespace"`
	MasterAddress     string `toml:"masterAddress"`
}

// Print prints the Config object to the log.
func (c Config) Print() {
	log.Printf("Vendor:\t\t%s\n", conf.Vendor)
	log.Printf("Version:\t\t%s\n", conf.Version)
	log.Printf("Executable:\t\t%s\n", conf.Exec)
	log.Printf("Database port:\t%s\n", conf.DBPort)
	log.Printf("Database addr:\t%s\n", conf.DBAddress)

	if conf.Vendor == "oracle" {
		log.Printf("SID:\t\t%s", conf.SID)
		log.Printf("Tablespace:\t\t%s", conf.DefaultTablespace)
	}

	log.Printf("Connector port:\t%s\n", conf.ConnectorPort)
	log.Printf("Username:\t\t%s\n", conf.User)
	log.Printf("Password:\t\t****\n")
	log.Printf("Master address\t%s\n", conf.MasterAddress)
}

// NewConfig returns a configuration file based on the vendor
func NewConfig(vendor string) Config {
	var conf Config

	switch vendor {
	case "mysql":
		conf = Config{
			Vendor:        "mysql",
			Version:       "5.5.53",
			DBPort:        "3306",
			DBAddress:     "localhost",
			ConnectorPort: "7000",
			User:          "root",
			Password:      "root",
			MasterAddress: "127.0.0.1",
		}

		switch runtime.GOOS {
		case "windows":
			conf.Exec = "C:\\path\\to\\mysql.exe"
		case "darwin":
			conf.Exec = "/usr/local/mysql/bin/mysql"
		default:
			conf.Exec = "/usr/bin/mysql"
		}
	case "postgres":
		conf = Config{
			Vendor:        "postgres",
			Version:       "9.4.9",
			DBPort:        "5432",
			DBAddress:     "localhost",
			ConnectorPort: "7000",
			User:          "postgres",
			Password:      "password",
			MasterAddress: "127.0.0.1",
		}

		switch runtime.GOOS {
		case "windows":
			conf.Exec = "C:\\path\\to\\psql.exe"
		case "darwin":
			conf.Exec = "/Library/PostgreSQL/9.4/bin/psql"
		default:
			conf.Exec = "/usr/bin/psql"
		}
	case "oracle":
		conf = Config{
			Vendor:            "oracle",
			Version:           "11g",
			DBPort:            "1521",
			DBAddress:         "localhost",
			SID:               "orcl",
			ConnectorPort:     "7000",
			User:              "system",
			Password:          "password",
			DefaultTablespace: "USERS",
			MasterAddress:     "127.0.0.1",
		}
		switch runtime.GOOS {
		case "windows":
			conf.Exec = "C:\\path\\to\\sqlplus.exe"
		case "darwin":
			conf.Exec = "/path/to/sqlplus"
		default:
			conf.Exec = "/path/to/sqlplus"
		}
	}

	return conf
}

func generateInteractive(filename string) (string, Config) {
	var (
		ok    = false
		vType = 0
	)

	for !ok {
		vType, ok = prompter.AskSelectionDef("What is the database vendor?", 0, vendors)
	}

	vendor := vendors[vType]
	def := NewConfig(vendor)

	var config Config

	config.Vendor = vendor
	config.Version = prompter.Ask("What is the database version?")
	config.DBPort = prompter.AskDef("What is the database port?", def.DBPort)
	config.DBAddress = prompter.AskDef("What is the database address?", def.DBAddress)

	if vendor == "oracle" {
		config.SID = prompter.AskDef("What is the SID?", def.SID)
		config.DefaultTablespace = prompter.AskDef("What is the default tablespace?", def.DefaultTablespace)
		config.Exec = prompter.AskDef("Where is the sqlplus executable?", def.Exec)
	} else if vendor == "mysql" {
		config.Exec = prompter.AskDef("Where is the mysql executable?", def.Exec)
	} else if vendor == "postgres" {
		config.Exec = prompter.AskDef("Where is the psql executable?", def.Exec)
	}

	config.User = prompter.AskDef("Who is the database user?", def.User)
	config.Password = prompter.AskDef("What is the database password?", def.Password)
	config.ConnectorPort = prompter.AskDef("What should the connector's port be?", def.ConnectorPort)
	config.MasterAddress = prompter.AskDef("What is the address of the Master server?", def.MasterAddress)

	fname := prompter.AskDef("What should we name the configuration file?", filename)

	return fname, config
}
