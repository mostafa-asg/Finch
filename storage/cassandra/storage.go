package cassandra

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gocql/gocql"
	config "github.com/spf13/viper"
)

const urls_table string = "urls"

type storage struct {
	session *gocql.Session
}

var ErrDuplicate = errors.New("Duplicate row ID")

func New() *storage {

	hosts := config.GetString("cassandra.hosts")
	cluster := gocql.NewCluster(hosts)
	cluster.Keyspace = config.GetString("cassandra.keyspace")
	cluster.Consistency = getConsistencyFromConfig()

	var session *gocql.Session
	var err error
	for {
		session, err = cluster.CreateSession()
		if err != nil {
			log.Println(fmt.Sprintf("Could not connect to casssandra at %s", hosts), err)
			//Try one second later
			time.Sleep(1 * time.Second)
		} else {
			break
		}
	}

	return &storage{
		session,
	}
}

func getConsistencyFromConfig() gocql.Consistency {

	consistency := strings.ToLower(config.GetString("cassandra.consistency"))
	switch consistency {
	case "any":
		return gocql.Any
	case "one":
		return gocql.One
	case "two":
		return gocql.Two
	case "three":
		return gocql.Three
	case "quorum":
		return gocql.Quorum
	case "all":
		return gocql.All
	case "localquorum":
		return gocql.LocalQuorum
	case "eachquorum":
		return gocql.EachQuorum
	case "localone":
		return gocql.LocalOne
	default:
		log.Fatal("invalid config value for cassandra.consistency")
	}

	//this line won't execute
	return gocql.One
}

func (st *storage) Put(id string, originalUrl string) error {
	insertStat := fmt.Sprintf("INSERT INTO %s(id,url) VALUES (?,?) IF NOT EXISTS;", urls_table)

	var (
		idCAS          string
		originalUrlCAS string
	)

	applied, err := st.session.Query(insertStat, id, originalUrl).ScanCAS(&idCAS, &originalUrlCAS)
	if applied == false {
		return ErrDuplicate
	}

	return err
}

func (st *storage) Get(id string) (string, error) {

	selectStat := fmt.Sprintf("SELECT url FROM %s WHERE id = ? LIMIT 1", urls_table)

	var url string
	err := st.session.Query(selectStat, id).Consistency(gocql.One).Scan(&url)
	if err != nil {
		return "", err
	}

	return url, nil
}
