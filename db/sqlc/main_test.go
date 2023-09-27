package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/okoroemeka/simple_bank/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries

var testDB *sql.DB

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
		return
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testQueries = New(testDB)
	os.Exit(m.Run())
}
