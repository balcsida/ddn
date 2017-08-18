package sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/djavorszky/ddn/server/database/data"
	"github.com/djavorszky/ddn/server/database/dbutil"
	"github.com/djavorszky/sutils"

	// DB
	_ "github.com/mattn/go-sqlite3"
)

// DB implements the BackendConnection
type DB struct {
	DBLocation string

	conn *sql.DB
}

// ConnectAndPrepare establishes a database connection and initializes the tables, if needed
func (lite *DB) ConnectAndPrepare() error {
	conn, err := sql.Open("sqlite3", lite.DBLocation)
	if err != nil {
		return fmt.Errorf("could not open connection: %v", err)
	}

	err = conn.Ping()
	if err != nil {
		return fmt.Errorf("database ping failed: %v", err)
	}
	lite.conn = conn

	err = lite.initTables()
	if err != nil {
		return fmt.Errorf("failed updating tables: %v", err)
	}

	return nil
}

// Close closes the database connection
func (lite *DB) Close() error {
	return lite.conn.Close()
}

// FetchByID returns the entry associated with that ID, or
// an error if it does not exist
func (lite *DB) FetchByID(ID int) (data.Row, error) {
	if err := lite.alive(); err != nil {
		return data.Row{}, fmt.Errorf("database down: %s", err.Error())
	}

	row := lite.conn.QueryRow("SELECT * FROM databases WHERE id = ?", ID)
	res, err := dbutil.ReadRow(row)
	if err != nil {
		return data.Row{}, fmt.Errorf("failed reading result: %v", err)
	}

	return res, nil
}

// FetchByCreator returns public entries that were created by the
// specified user, an empty list if it's not the user does
// not have any entries, or an error if something went
// wrong
func (lite *DB) FetchByCreator(creator string) ([]data.Row, error) {
	return []data.Row{}, nil
}

// FetchPublic returns all entries that have "Public" set to true
func (lite *DB) FetchPublic() ([]data.Row, error) {
	return []data.Row{}, nil
}

// FetchAll returns all entries.
func (lite *DB) FetchAll() ([]data.Row, error) {
	return []data.Row{}, nil
}

// Insert adds an entry to the database, returning its ID
func (lite *DB) Insert(row *data.Row) error {
	if err := lite.alive(); err != nil {
		return fmt.Errorf("database down: %s", err.Error())
	}

	query := "INSERT INTO `databases` (`dbname`, `dbuser`, `dbpass`, `dbsid`, `dumpfile`, `createDate`, `expiryDate`, `creator`, `connectorName`, `dbAddress`, `dbPort`, `dbvendor`, `status`, `message`, `visibility`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?,?, ?, ?, ?, ?)"

	res, err := lite.conn.Exec(query,
		row.DBName,
		row.DBUser,
		row.DBPass,
		row.DBSID,
		row.Dumpfile,
		row.CreateDate,
		row.ExpiryDate,
		row.Creator,
		row.ConnectorName,
		row.DBAddress,
		row.DBPort,
		row.DBVendor,
		row.Status,
		row.Message,
		row.Public,
	)
	if err != nil {
		return fmt.Errorf("insert failed: %v", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed getting new ID: %v", err)
	}

	row.ID = int(id)

	return nil
}

// Update updates an already existing entry
func (lite *DB) Update(row *data.Row) error {
	return nil
}

// Delete removes the entry from the database
func (lite *DB) Delete(row data.Row) error {
	return nil
}

type dbUpdate struct {
	Query   string
	Comment string
}

var queries = []dbUpdate{
	dbUpdate{
		Query:   "CREATE TABLE version (queryId INTEGER PRIMARY KEY, query TEXT NULL, comment TEXT NULL, date DATETIME NULL);",
		Comment: "Create the version table",
	},
	dbUpdate{
		Query:   "CREATE TABLE databases (id INTEGER PRIMARY KEY, dbname VARCHAR(255) NULL, dbuser VARCHAR(255) NULL, dbpass VARCHAR(255) NULL, dbsid VARCHAR(45) NULL, dumpfile TEXT NULL, createDate DATETIME NULL, expiryDate DATETIME NULL, creator VARCHAR(255) NULL, connectorName VARCHAR(255) NULL, dbAddress VARCHAR(255) NULL, dbPort VARCHAR(45) NULL, dbvendor VARCHAR(255) NULL, status INTEGER, message TEXT, visibility INTEGER DEFAULT 0);",
		Comment: "Create the databases table",
	},
	dbUpdate{
		Query:   "UPDATE databases SET message = '' WHERE message IS NULL;",
		Comment: "Update 'message' columns to empty where null",
	},
}

func (lite *DB) initTables() error {
	var (
		err      error
		startLoc int
	)

	lite.conn.QueryRow("SELECT count(*) FROM version").Scan(&startLoc)

	for _, q := range queries[startLoc:] {
		log.Printf("Updating database %q", q.Comment)
		_, err = lite.conn.Exec(q.Query)
		if err != nil {
			return fmt.Errorf("executing query %q (%q) failed: %s", q.Comment, q.Query, sutils.TrimNL(err.Error()))
		}

		_, err = lite.conn.Exec("INSERT INTO version (query, comment, date) VALUES (?, ?, ?)", q.Query, q.Comment, time.Now())
		if err != nil {
			return fmt.Errorf("updating version table with query %q (%q) failed: %s", q.Comment, q.Query, sutils.TrimNL(err.Error()))
		}
	}

	return nil
}

// Alive checks whether the connection is alive. Returns error if not.
func (lite *DB) alive() error {
	defer func() {
		if p := recover(); p != nil {
			log.Println("Panic Attack! Database seems to be down.")
		}
	}()

	_, err := lite.conn.Exec("select * from databases WHERE 1 = 0")
	if err != nil {
		return fmt.Errorf("executing stayalive query failed: %s", sutils.TrimNL(err.Error()))
	}

	return nil
}
