package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	config "github.com/spf13/viper"
)

const urls_table string = "urls"

type storage struct {
	dbPath string
}

func New() *storage {
	err := ensureDatabaseExists()
	if err != nil {
		log.Fatal("Error in initializing sqllite database", err)
	}
	return &storage{
		config.GetString("sqlite.path"),
	}
}

func (st *storage) Put(id string, originalUrl string) error {

	db, err := sql.Open("sqlite3", st.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("insert into %s(id, url) values('%s','%s')", urls_table, id, originalUrl))
	return err
}

func (st *storage) Get(id string) (string, error) {

	db, err := sql.Open("sqlite3", st.dbPath)
	if err != nil {
		return "", err
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("select url from %s where id='%s'", urls_table, id))
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

	sqlStmt := fmt.Sprintf("create table if not exists %s (id nvarchar(10) not null primary key, url text);", urls_table)

	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
