package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	config "github.com/spf13/viper"
)

type storage struct {
	dbPath string
	table  string
}

func New() *storage {
	ensureDatabaseExists()
	return &storage{
		config.GetString("sqlite.path"),
		config.GetString("sqlite.table"),
	}
}

func (st *storage) Put(id string, originalUrl string) error {

	db, err := sql.Open("sqlite3", st.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("insert into %s(id, url) values('%s','%s')", st.table, id, originalUrl))
	return err
}

func (st *storage) Get(id string) (string, error) {

	db, err := sql.Open("sqlite3", st.dbPath)
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

	db, err := sql.Open("sqlite3", config.GetString("sqlite.path"))
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStmt := fmt.Sprintf("create table if not exists %s (id nvarchar(10) not null primary key, url text);", config.GetString("sqlite.table"))

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
