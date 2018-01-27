package users

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/mostafa-asg/finch/generator/base62"
	"github.com/mostafa-asg/finch/http/model"
	"github.com/mostafa-asg/finch/test/cache"
)

type user struct {
	servers       []string
	currentServer int
	done          chan bool
	wg            *sync.WaitGroup
}

func newUser(servers []string, waitGroup *sync.WaitGroup) *user {

	size := len(servers)
	if size == 0 {
		log.Fatal("No servers to send requests")
	}

	result := &user{
		servers: make([]string, size),
		done:    make(chan bool),
		wg:      waitGroup,
	}
	copy(result.servers, servers)
	return result
}

func (u *user) makeGetRequests() {

	go func() {

		u.wg.Add(1)
		ticker := time.NewTicker(100 * time.Millisecond)

		for {

			select {
			case <-ticker.C:
				tinyURL, originalURL := cache.GetInstance().ReadRandom()
				res, err := http.Get(fmt.Sprintf("%s/get/%s", u.nextServer(), tinyURL))

				if err != nil {
					log.Fatal(err.Error())
				}

				var response model.GetResponse
				json.NewDecoder(res.Body).Decode(&response)

				if response.Found == false {
					log.Fatalf("tiny url %s not found", tinyURL)
				}

				if response.Url != originalURL {
					log.Fatalf("Invalid url. Expected %s but found %s for tiny url %s", originalURL, response.Url, tinyURL)
				}

				res.Body.Close()
			case <-u.done:
				ticker.Stop()
				wg.Done()
				return
			} //end select

		} //end for

	}()
}

func (u *user) makeHashRequests() {

	go func() {

		generator := base62.New(10)
		ticker := time.NewTicker(100 * time.Millisecond)
		u.wg.Add(1)

		for {

			select {
			case <-ticker.C:
				body := `
				{
					"url":"%s"
				}
				`
				url := generator.GenerateID()
				body = fmt.Sprintf(body, url)
				res, err := http.Post(fmt.Sprintf("%s/hash", u.nextServer()), "application/json", strings.NewReader(body))
				if err != nil {
					log.Fatal(err.Error())
				}

				var response model.HashResponse
				json.NewDecoder(res.Body).Decode(&response)

				if response.ErrorMessage != "" {
					log.Fatal(response.ErrorMessage)
				} else {
					cache.GetInstance().Write(response.TinyUrl, url)
				}

				res.Body.Close()

			case <-u.done:
				ticker.Stop()
				u.wg.Done()
				return
			} //end select

		} //end for
	}()

}

func (u *user) nextServer() string {

	size := len(u.servers)
	u.currentServer = u.currentServer + 1
	return u.servers[u.currentServer%size]

}

func (u *user) stop() {
	close(u.done)
}
