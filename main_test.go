package main

import (
	"fmt"
	"math/rand"
	"path"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestAddCredential(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "maximum effort",
	}

	dbName := fmt.Sprintf("credential_testing_%d", rand.Int())

	dbPath := path.Join("test_dbs", dbName)

	db, _ := OpenSqliteConn(dbPath)
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
