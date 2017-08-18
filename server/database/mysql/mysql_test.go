package mysql

import (
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/djavorszky/ddn/server/database/data"
	"github.com/djavorszky/ddn/server/database/dbutil"
	"github.com/djavorszky/sutils"

	_ "github.com/go-sql-driver/mysql"
)

const (
	testAddr = "127.0.0.1"
	testPort = "3306"
	testUser = "travis"
	testPass = ""
	testName = "unit_test"
)

var (
	testConn *sql.DB

	mys DB

	gmt, _ = time.LoadLocation("GMT")

	testEntry = data.Row{
		ID:            1,
		DBName:        "testDB",
		DBUser:        "testUser",
		DBPass:        "testPass",
		DBSID:         "testsid",
		Dumpfile:      "testloc",
		CreateDate:    time.Now().In(gmt),
		ExpiryDate:    time.Now().In(gmt).AddDate(0, 0, 30),
		Creator:       "test@gmail.com",
		ConnectorName: "mysql-55",
		DBAddress:     "localhost",
		DBPort:        "3306",
		DBVendor:      "mysql",
		Message:       "",
		Status:        100,
	}
)

func TestMain(m *testing.M) {
	// For these tests to run, a local database should be present
	// which has a user named 'travis' with no password authentication

	err := setup()
	if err != nil {
		fmt.Printf("Failed setup: %s", err.Error())
		os.Exit(-1)
	}

	res := m.Run()

	err = teardown()
	if err != nil {
		fmt.Printf("Failed teardown: %s", err.Error())
		os.Exit(-1)
	}

	os.Exit(res)
}

// Connect to a local database
func setup() error {
	var err error

	datasource := fmt.Sprintf("%s:%s@tcp(%s:%s)/", testUser, testPass, testAddr, testPort)
	err = testConnDS(datasource)
	if err != nil {
		return fmt.Errorf("failed to setup test connection: %v", err)
	}

	_, err = testConn.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARSET utf8;", testName))
	if err != nil {
		return fmt.Errorf("failed creating database: %s", sutils.TrimNL(err.Error()))
	}

	testConn.Close()

	err = testConnDS(datasource + testName)
	if err != nil {
		return fmt.Errorf("failed connecting to created database")
	}

	mainDS := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", testUser, testPass, testAddr, testPort, testName)
	err = mys.connect(mainDS)
	if err != nil {
		return fmt.Errorf("failed initializing main connection")
	}

	return nil
}

// DROP EVERYTHING!!4one
func teardown() error {
	var err error

	_, err = testConn.Exec(fmt.Sprintf("DROP DATABASE %s;", testName))
	if err != nil {
		return fmt.Errorf("failed dropping database: %s", sutils.TrimNL(err.Error()))
	}

	testConn.Close()
	mys.Close()

	return nil
}

func testConnDS(datasource string) error {
	var err error

	testConn, err = sql.Open("mysql", datasource+"?parseTime=true")
	if err != nil {
		return fmt.Errorf("creating connection pool failed: %s", err.Error())
	}

	err = testConn.Ping()
	if err != nil {
		testConn.Close()
		return fmt.Errorf("database ping failed: %s", sutils.TrimNL(err.Error()))
	}

	return nil
}

func TestInitTables(t *testing.T) {
	var err error

	_, err = testConn.Exec("SELECT 1 FROM version LIMIT 1;")
	if err == nil {
		t.Errorf("Version table already exists before test even ran.")
	}

	_, err = testConn.Exec("SELECT 1 FROM `databases` LIMIT 1;")
	if err == nil {
		t.Errorf("Databases table already exists before test even ran.")
	}

	err = mys.initTables()
	if err != nil {
		t.Errorf("Failed initializing tables: %s", err.Error())
	}

	_, err = testConn.Exec("SELECT 1 FROM version LIMIT 1;")
	if err != nil {
		t.Errorf("Version table has not been created.")
	}

	_, err = testConn.Exec("SELECT 1 FROM `databases` LIMIT 1;")
	if err != nil {
		t.Errorf("Databases table has not been created.")
	}

	type versiontest struct {
		queryID int
		query   string
		comment string
		date    time.Time
	}

	rows, _ := testConn.Query("SELECT * FROM version")

	for rows.Next() {
		var row versiontest

		rows.Scan(&row.queryID, &row.query, &row.comment, &row.date)

		dbu := queries[row.queryID-1]

		if row.query != dbu.Query {
			t.Errorf("Saved query not what was expected")
		}

		if row.comment != dbu.Comment {
			t.Errorf("Saved comment not what was expected")
		}
	}
	err = rows.Err()
	if err != nil {
		t.Errorf("error reading result from query: %s", err.Error())
	}
}

func TestFetchByID(t *testing.T) {
	mys.Insert(&testEntry)

	res, err := mys.FetchByID(testEntry.ID)
	if err != nil {
		t.Errorf("FetchById(%d) failed with error: %v", testEntry.ID, err)
	}

	if err := dbutil.CompareRows(res, testEntry); err != nil {
		t.Errorf("Fetched result not the same as queried: %v", err)
	}
}

func TestFetchByCreator(t *testing.T) {
	creator := "someone@somewhere.com"

	testEntry.Creator = creator

	mys.Insert(&testEntry)
	mys.Insert(&testEntry)

	results, err := mys.FetchByCreator(creator)
	if err != nil {
		t.Errorf("failed to fetch by creator: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected resultset to have 2 results, %d instead", len(results))
	}

	for _, res := range results {
		if res.Creator != creator {
			t.Errorf("Creator mismatch: Got %q, expected %q", res.Creator, creator)
		}
	}
}

func TestInsert(t *testing.T) {
	err := mys.Insert(&testEntry)
	if err != nil {
		t.Errorf("mys.Insert(testEntry) failed with error: %v", err)
	}

	if testEntry.ID == 0 {
		t.Errorf("mys.Insert(testEntry) resulted in id of 0")
	}

	result, err := mys.FetchByID(testEntry.ID)
	if err != nil {
		t.Errorf("FetchById(%d) resulted in error: %v", testEntry.ID, err)
	}

	if err = dbutil.CompareRows(testEntry, result); err != nil {
		t.Errorf("Persisted and read results not the same: %v", err)
	}
}

func TestUpdate(t *testing.T) {
	mys.Insert(&testEntry)

	// We're updating by ID - this should updated the row for "testEntry"
	updatedEntry := data.Row{
		ID:            testEntry.ID,
		DBName:        "updatedtestDB",
		DBUser:        "updatedtestUser",
		DBPass:        "updatedtestPass",
		DBSID:         "updatedtestsid",
		Dumpfile:      "updatedtestloc",
		CreateDate:    time.Now().In(gmt),
		ExpiryDate:    time.Now().In(gmt).AddDate(0, 0, 30),
		Creator:       "updatedtest@gmail.com",
		ConnectorName: "updatedysql-55",
		DBAddress:     "updatedlocalhost",
		DBPort:        "updated3306",
		DBVendor:      "updatedmysql",
		Message:       "updated",
		Status:        200,
	}

	err := mys.Update(&updatedEntry)
	if err != nil {
		t.Errorf("Update(updatedEntry) failed: %v", err)
	}

	readEntry, _ := mys.FetchByID(testEntry.ID)

	if err := dbutil.CompareRows(updatedEntry, readEntry); err != nil {
		t.Errorf("Updated and read entries not the same: %v", err)
	}
}

func TestDelete(t *testing.T) {
	mys.Insert(&testEntry)

	err := mys.Delete(testEntry)
	if err != nil {
		t.Errorf("Delete failed: %v", err)
	}

	_, err = mys.FetchByID(testEntry.ID)
	if err == nil {
		t.Errorf("Row not deleted")
	}
}

func TestFetchPublic(t *testing.T) {
	res, err := mys.FetchPublic()
	if err != nil {
		t.Errorf("FetchPublic() error: %v", err)
	}

	if len(res) != 0 {
		t.Errorf("FetchPublic() returned with entries, shouldn't have")
	}

	testEntry.Public = 1

	mys.Insert(&testEntry)

	res, err = mys.FetchPublic()
	if err != nil {
		t.Errorf("FetchPublic() error: %v", err)
	}

	if len(res) != 1 {
		t.Errorf("FetchPublic() expected 1 result, got %d instead", len(res))
	}

	if err := dbutil.CompareRows(res[0], testEntry); err != nil {
		t.Errorf("Read and persisted mismatch: %v", err)
	}
}

func TestFetchAll(t *testing.T) {
	var count int

	mys.conn.QueryRow("SELECT count(*) FROM `databases`").Scan(&count)

	entries, err := mys.FetchAll()
	if err != nil {
		t.Errorf("FetchAll() encountered error: %v", err)
	}

	if len(entries) != count {
		t.Errorf("Expected size %d, got %d instead", count, len(entries))
	}
}

func TestReadRow(t *testing.T) {
	err := mys.Insert(&testEntry)
	if err != nil {
		t.Errorf("Failed adding a entry: %s", err.Error())
	}

	rows, err := testConn.Query("SELECT * FROM `databases` WHERE id = ?", testEntry.ID)
	if err != nil {
		t.Errorf("Failed querying for entries: %s", err.Error())
	}

	for rows.Next() {
		row, err := dbutil.ReadRows(rows)
		if err != nil {
			t.Errorf("Failed reading row from rows: %s", err.Error())
		}

		if err = dbutil.CompareRows(testEntry, row); err != nil {
			t.Errorf("Persisted and read DBEntry not the same: %s", err.Error())
		}
	}

	// cleanup
	_, err = testConn.Exec("DELETE FROM `databases` WHERE ID = ?", testEntry.ID)
	if err != nil {
		t.Errorf("Could not delete created entry")
	}

	testEntry.ID++
}
