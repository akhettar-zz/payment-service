package api_test

import (
	"payment-service/repository"
	"io/ioutil"
	"os"
	"testing"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/dbtest"
)

const (
	DBName = "test"
	dir    = "testing-db"
)

var Server dbtest.DBServer
var Session *mgo.Session

var Repository *repository.MongoRepository

// TestMain wraps all tests with the needed initialized mock DB and fixtures
// This test runs before other integration test. It starts an instance of mongo db in the background (provided you have mongo
// installed on the server on which this test will be running) and shuts it down.
func TestMain(m *testing.M) {

	// The tempdir is created so MongoDB has a location to store its files.
	// Contents are wiped once the server stops
	tempDir, _ := ioutil.TempDir("", dir)
	Server.SetPath(tempDir)

	// My main session var is now set to the temporary MongoDB instance
	Session = Server.Session()

	Repository = &repository.MongoRepository{Session}

	// Run the test suite
	retCode := m.Run()

	// Make sure we DropDatabase so we make absolutely sure nothing is left or locked while wiping the data and
	// close session
	Session.DB(DBName).DropDatabase()
	Session.Close()

	// Stop shuts down the temporary server and removes data on disk.
	Server.Stop()

	os.RemoveAll(tempDir)

	// call with result of m.Run()
	os.Exit(retCode)
}
