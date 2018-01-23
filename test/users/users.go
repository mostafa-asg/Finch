package users

import (
	"sync"
)

var usersList []*user
var wg *sync.WaitGroup

func MakeRequests(numberOfhashUsers int, numberOfGetUsers int, servers []string) {

	usersList = make([]*user, numberOfhashUsers+numberOfGetUsers)
	wg = new(sync.WaitGroup)

	index := 0
	for i := 1; i <= numberOfhashUsers; i++ {
		usersList[index] = newUser(servers, wg)
		usersList[index].makeHashRequests()
		index++
	}

	for i := 1; i <= numberOfGetUsers; i++ {
		usersList[index] = newUser(servers, wg)
		usersList[index].makeGetRequests()
		index++
	}
}

func Stop() {

	for _, user := range usersList {
		user.stop()
	}
	wg.Wait()

}
