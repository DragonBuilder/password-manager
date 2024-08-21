package main

import (
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

	db, _ := OpenSqliteConn("tests/db")
	defer db.Close()
	service := Service{db}

	err := service.AddCredential(cred)
	assert.NoError(t, err)

	row := db.QueryRow("SELECT username FROM credentials WHERE username='deadpool@gmail.com'")
	got := Credential{}

	err = row.Scan(&got.Username)
	assert.NoError(t, err)

	assert.Equal(t, cred.Username, got.Username)
}
