package mysql

import (
	"database/sql"
	"errors"
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
		"insert into %s(shortUrl,ts,referrer,browser,country,os) "+
			"values ('%s','%s','%s','%s','%s','%s')",
		visit_table, shortUrl, info.Time.Format("2006-01-02 15:04:05"), info.Referrer, info.Browser, info.Country, info.OS)

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

func (st *storage) GetStats(shortUrl string, queryType string) (core.Stats, error) {
	switch queryType {
	case "all":
		return st.getAllTimeStats(shortUrl)
	case "month":
		return st.getLastMonthStats(shortUrl)
	case "week":
		return st.getLastWeekStats(shortUrl)
	case "day":
		return st.getTodayStats(shortUrl)
	default:
		return core.Stats{}, errors.New("Invalid query type")
	}
}

func getQuery(groupBy, shortUrl, endTime string) string {

	var timeClouse string

	if endTime == "all" {
		timeClouse = ""
	} else {
		timeClouse = fmt.Sprintf("AND ts>=DATE_SUB(now() , INTERVAL %s DAY)", endTime)
	}

	r := fmt.Sprintf("SELECT %s , count(*) as clicks FROM ("+
		"SELECT * from visit where shortUrl='%s' AND ts<=now() %s)"+
		" AS t1 GROUP BY %s;", groupBy, shortUrl, timeClouse, groupBy)
	println(r)
	return r
}

func (st *storage) executeQuery(shortUrl, endTime string) (core.Stats, error) {
	timelines, err := st.execute(getQuery("YEAR(ts)", shortUrl, endTime))

	if err != nil {
		return core.Stats{}, nil
	}

	referrals, err := st.execute(getQuery("referrer", shortUrl, endTime))
	if err != nil {
		return core.Stats{}, nil
	}

	browsers, err := st.execute(getQuery("browser", shortUrl, endTime))
	if err != nil {
		return core.Stats{}, nil
	}

	countries, err := st.execute(getQuery("country", shortUrl, endTime))
	if err != nil {
		return core.Stats{}, nil
	}

	opertingSystems, _ := st.execute(getQuery("os", shortUrl, endTime))
	if err != nil {
		return core.Stats{}, nil
	}

	return core.Stats{
		Timeline:  timelines,
		Referrals: referrals,
		Browsers:  browsers,
		Countries: countries,
		OS:        opertingSystems,
	}, nil
}

func (st *storage) getAllTimeStats(shortUrl string) (core.Stats, error) {
	return st.executeQuery(shortUrl, "all")
}

func (st *storage) getLastMonthStats(shortUrl string) (core.Stats, error) {
	return st.executeQuery(shortUrl, "30")
}

func (st *storage) getLastWeekStats(shortUrl string) (core.Stats, error) {
	return st.executeQuery(shortUrl, "7")
}

func (st *storage) getTodayStats(shortUrl string) (core.Stats, error) {
	return st.executeQuery(shortUrl, "1")
}

func (st *storage) execute(query string) (map[string]int, error) {
	result := make(map[string]int)

	db, err := sql.Open("mysql", st.connectionStr)
	if err != nil {
		return map[string]int{}, err
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		return map[string]int{}, err
	}

	var key string
	var count int

	for rows.Next() {
		err = rows.Scan(&key, &count)
		if err != nil {
			return result, err
		}
		result[key] = count
	}

	return result, nil
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
		"ts timestamp ,"+
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
