package main

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
	_ "modernc.org/sqlite"
)

type Credential struct {
	Website  string
	Username string
	Password string
}

func main() {
	// (&cli.App{}).Run(os.Args)

	app := &cli.App{
		Name:  "pman",
		Usage: "manage your passwords",
	}

	app.Run(os.Args)
}

func OpenSqliteConn(dbPath string) (*sql.DB, error) {
	dir, _ := path.Split(dbPath)
	if dir != "" {
		os.MkdirAll(dir, os.ModePerm)
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("error opening connection to sqlite3 : %w", err)
	}
	return db, nil
}

type Service struct {
	db *sql.DB
}

func (s *Service) AddCredential(cred Credential) error {
	const create_credentials_table_qry = `
CREATE TABLE IF NOT EXISTS credentials (
	id INTEGER NOT NULL PRIMARY KEY,
	website TEXT NOT NULL,
	username TEXT NOT NULL,
	password TEXT NOT NULL,
	created_at DATETIME NOT NULL,
	updated_at DATETIME NOT NULL,
	UNIQUE(website, username)
);
`

	if _, err := s.db.Exec(create_credentials_table_qry); err != nil {
		return fmt.Errorf("error running query to create credentials table : %w", err)
	}

	const insert_credential_qry = `
INSERT INTO credentials VALUES(NULL, ?, ?, ?, ?, ?)
`

	now := time.Now()
	if _, err := s.db.Exec(insert_credential_qry, cred.Website, cred.Username, cred.Password, now, now); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return errors.New("CredentialsAlreadyExistsForWebsite")
		}
		return fmt.Errorf("error inserting credential to db: %w", err)
	}

	return nil
}

func (s *Service) RetrieveCredential(website string) ([]Credential, error) {
	const get_credential_qry = `
SELECT username, password FROM credentials WHERE website=?
`
	rows, err := s.db.Query(get_credential_qry, website)
	if err != nil {
		return nil, fmt.Errorf("error fetching credentials : %w", err)
	}
	var results = make([]Credential, 0)

	for rows.Next() {
		var cred = Credential{}
		rows.Scan(&cred.Username, &cred.Password)
		results = append(results, cred)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("credentials not found")
	}

	return results, nil

}
