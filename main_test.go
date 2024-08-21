package main

import (
	"fmt"
	"math/rand"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	_ "modernc.org/sqlite"
)

func dbPathForTest(dbName string) string {
	dbName = fmt.Sprintf("%s_%d", dbName, rand.Int())
	dbPath := path.Join("test_dbs", dbName)
	return dbPath
}

func TestOpenSqliteConn(t *testing.T) {
	db, err := OpenSqliteConn(dbPathForTest("conn_test"))
	assert.NoError(t, err)
	defer db.Close()
	assert.NotNil(t, db)
}

func TestAddCredential(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "maximum effort",
	}

	db, _ := OpenSqliteConn(dbPathForTest("test_add_credential"))
	defer db.Close()
	service := Service{db}

	err := service.AddCredential(cred)
	assert.NoError(t, err)

	row := db.QueryRow("SELECT username FROM credentials WHERE username=?", cred.Username)
	got := Credential{}

	err = row.Scan(&got.Username)
	assert.NoError(t, err)

	assert.Equal(t, cred.Username, got.Username)
}
