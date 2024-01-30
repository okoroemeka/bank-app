package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/okoroemeka/simple_bank/util"
	"log"
	"os"
	"testing"
)

var testStore Store

func TestMain(m *testing.M) {
	config, err := util.LoadConfig("../..")
	if err != nil {
		log.Fatal("cannot load config", err)
		return
	}

	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	testStore = NewStore(connPool)

	os.Exit(m.Run())
}
