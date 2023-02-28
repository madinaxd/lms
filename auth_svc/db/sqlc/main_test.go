package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	driver     = "postgres"
	dataSource = "postgresql://postgres:postgres@localhost:5432/auth_svc?sslmode=disable"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	var err error
	testDB, err = sql.Open(driver, dataSource)
	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}
	testQueries = New(testDB)

	os.Exit(m.Run())
}
