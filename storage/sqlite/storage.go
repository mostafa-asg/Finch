package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/mostafa-asg/finch/core"
	config "github.com/spf13/viper"
)

const urls_table string = "urls"
const visit_table string = "visit"

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

func (st *storage) Visit(shortUrl string, info core.VisitInfo) error {
	db, err := sql.Open("sqlite3", st.dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	sqlStm := fmt.Sprintf(
		"insert into %s(shortUrl,year,month,day,hour,minute,referrer,browser,country,os) "+
			"values ('%s',%d,%d,%d,%d,%d,'%s','%s','%s','%s')",
		visit_table, shortUrl, info.Year, info.Month, info.Day, info.Hour,
		info.Minute, info.Referrer, info.Browser, info.Country, info.OS)

	_, err = db.Exec(sqlStm)
	return err
}

func (st *storage) Count(shortUrl string) int {
	db, err := sql.Open("sqlite3", st.dbPath)
	if err != nil {
		log.Println("Could not open sqlite database", err)
		return 0
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("select count(*) from %s where shortUrl='%s'", visit_table, shortUrl))
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Panicln("Sqlite count error", err)
		return 0
	}

	return count
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

	sqlStmt = fmt.Sprintf("create table if not exists %s ("+
		"id integer primary key, "+
		"shortUrl text ,"+
		"year integer ,"+
		"month integer ,"+
		"day integer ,"+
		"hour integer ,"+
		"minute integer ,"+
		"referrer text ,"+
		"browser text ,"+
		"country text ,"+
		"os text);", visit_table)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return err
	}

	return nil
}
