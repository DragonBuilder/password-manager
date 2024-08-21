package main

import (
	"fmt"
	"math/rand"
	"os"
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
	dbPath := dbPathForTest("conn_test")
	db, err := OpenSqliteConn(dbPath)

	assert.NoError(t, err)

	defer os.Remove(dbPath)
	defer db.Close()

	assert.NotNil(t, db)
	assert.NoError(t, db.Ping())
}

func TestAddCredential(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "maximum effort",
	}

	dbPath := dbPathForTest("test_add_credential")
	db, _ := OpenSqliteConn(dbPath)

	defer os.Remove(dbPath)
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

func TestRetrieveCredentialShouldReturnCorrectCredentialWhenCredentialExistsForAWebsite(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "wolverine@gmail.com",
		Password: "no",
	}

	dbPath := dbPathForTest("test_retrieve_credential")
	db, _ := OpenSqliteConn(dbPath)

	defer os.Remove(dbPath)
	defer db.Close()

	srv := Service{db}

	err := srv.AddCredential(cred)
	assert.NoError(t, err)

	got_creds, err := srv.RetrieveCredential(cred.Website)
	assert.NoError(t, err)

	assert.Equal(t, 1, len(got_creds))

	assert.Equal(t, cred.Username, got_creds[0].Username)
	assert.Equal(t, cred.Password, got_creds[0].Password)
}

func TestAddMultipleCredentialsShouldAllowNewCredentialForSameWebsiteWhenUsernameIsDifferent(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "maximum effort",
	}

	dbPath := dbPathForTest("test_add_credential")
	db, _ := OpenSqliteConn(dbPath)

	defer os.Remove(dbPath)
	defer db.Close()

	srv := Service{db}

	err := srv.AddCredential(cred)
	assert.NoError(t, err)

	err = srv.AddCredential(Credential{
		Website:  "gmail",
		Username: "funpool@gmail.com",
		Password: "minimum effort",
	})
	assert.NoError(t, err)

	got_creds, err := srv.RetrieveCredential(cred.Website)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(got_creds))
}
func TestAddMultipleCredentialsShouldDisAllowNewCredentialForSameWebsiteWhenUsernameIsAlreadyExists(t *testing.T) {
	cred := Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "maximum effort",
	}

	dbPath := dbPathForTest("test_add_credential")
	db, _ := OpenSqliteConn(dbPath)

	defer os.Remove(dbPath)
	defer db.Close()

	srv := Service{db}

	err := srv.AddCredential(cred)
	assert.NoError(t, err)

	err = srv.AddCredential(Credential{
		Website:  "gmail",
		Username: "deadpool@gmail.com",
		Password: "minimum effort",
	})
	assert.Error(t, err)
}
