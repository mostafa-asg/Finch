package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

const dbPath = "./finch.db"
const table = "urls"

type storage struct {
}

func New() *storage {
	ensureDatabaseExists()
	return &storage{}
}

func (st *storage) Put(id string, originalUrl string) error {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_ , err = db.Exec(fmt.Sprintf("insert into %s(id, url) values('%s','%s')", table, id, originalUrl))
	return err
}

func (st *storage) Get(id string) (string, error) {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("select url from urls where id='%s'", id))
	var url string
	err = row.Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}

func ensureDatabaseExists() error {

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := fmt.Sprintf("create table if not exists %s (id nvarchar(10) not null primary key, url text);", table)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
