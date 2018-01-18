package cassandra

import (
	"fmt"

	"github.com/gocql/gocql"
)

type storage struct {
	session *gocql.Session
	table   string
}

func New() *storage {

	cluster := gocql.NewCluster("127.0.0.1")
	cluster.Keyspace = "default"
	cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	if err != nil {
		panic(err)
	}

	return &storage{
		session,
		"urls",
	}
}

func (st *storage) Put(id string, originalUrl string) error {

	insertStat := fmt.Sprintf("INSERT INTO %s(id,url) VALUES (?,?) IF NOT EXISTS;", st.table)
	return st.session.Query(insertStat, id, originalUrl).Exec()
}

func (st *storage) Get(id string) (string, error) {

	selectStat := fmt.Sprintf("SELECT url FROM %s WHERE id = ? LIMIT 1", st.table)

	var url string
	err := st.session.Query(selectStat, id).Consistency(gocql.One).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}
