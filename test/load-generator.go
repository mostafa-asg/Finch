package main

import (
	"flag"
	"strings"
	"time"

	"github.com/mostafa-asg/finch/test/users"
)

func main() {

	var servers string
	var hashUsers int
	var getUsers int
	var duration int64
	flag.StringVar(&servers, "servers", "http://localhost:8585", "Comma seperated list of finch servers")
	flag.IntVar(&hashUsers, "writeUsers", 5, "Number of users that insert a record to database")
	flag.IntVar(&getUsers, "readUsers", 15, "Number of users that read a record from database")
	flag.Int64Var(&duration, "duration", 10, "How many seconds does it take?")
	flag.Parse()

	go func() {
		users.MakeRequests(hashUsers, getUsers, strings.Split(servers, ","))
	}()

	time.Sleep(time.Duration(duration) * time.Second)

	users.Stop()

}
