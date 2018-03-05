package mysql

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/mostafa-asg/finch/core"
	config "github.com/spf13/viper"
)

const urls_table string = "urls"
const visit_table string = "visit"

type storage struct {
	connectionStr string
}

func New() *storage {
	cs, err := ensureDatabaseExists()
	if err != nil {
		log.Fatal("Error in initializing mysql database", err)
	}
	return &storage{
		connectionStr: cs,
	}
}

func (st *storage) Put(id string, originalUrl string) error {

	db, err := sql.Open("mysql", st.connectionStr)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("insert into %s(id, url) values('%s','%s')", urls_table, id, originalUrl))
	return err
}

func (st *storage) Get(id string) (string, error) {

	db, err := sql.Open("mysql", st.connectionStr)
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
	db, err := sql.Open("mysql", st.connectionStr)
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
	db, err := sql.Open("mysql", st.connectionStr)
	if err != nil {
		log.Println("Could not open mysql database", err)
		return 0
	}
	defer db.Close()

	row := db.QueryRow(fmt.Sprintf("select count(*) from %s where shortUrl='%s'", visit_table, shortUrl))
	var count int
	err = row.Scan(&count)
	if err != nil {
		log.Panicln("mysql count error", err)
		return 0
	}

	return count
}

func ensureDatabaseExists() (string, error) {

	connectionStr := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		config.GetString("mysql.user"),
		config.GetString("mysql.pass"),
		config.GetString("mysql.host"),
		config.GetString("mysql.database"))

	db, err := sql.Open("mysql", connectionStr)
	if err != nil {
		return "", err
	}
	defer db.Close()

	sqlStmt := fmt.Sprintf("create table if not exists %s (id nvarchar(10) not null primary key, url text);", urls_table)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return "", err
	}

	sqlStmt = fmt.Sprintf("create table if not exists %s ("+
		"id integer AUTO_INCREMENT, "+
		"shortUrl text ,"+
		"year integer ,"+
		"month integer ,"+
		"day integer ,"+
		"hour integer ,"+
		"minute integer ,"+
		"referrer text ,"+
		"browser text ,"+
		"country text ,"+
		"os text , PRIMARY KEY (id));", visit_table)
	_, err = db.Exec(sqlStmt)
	if err != nil {
		return "", err
	}

	return connectionStr, nil
}
