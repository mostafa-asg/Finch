package cache

import (
	"math/rand"
	"sync"
	"time"
)

type tinyUrl struct {
	tiny        string
	originalUrl string
}

type readRequest struct {
	response chan tinyUrl
}

type cache struct {
	done        chan bool
	writeBuffer chan tinyUrl
	readBuffer  chan readRequest
	urls        []tinyUrl
	rnd         *rand.Rand
}

var instance *cache
var lock sync.Mutex

func GetInstance() *cache {

	if instance != nil {
		return instance
	}

	lock.Lock()
	defer lock.Unlock()
	if instance == nil {

		instance = &cache{
			done:        make(chan bool),
			writeBuffer: make(chan tinyUrl, 10),
			readBuffer:  make(chan readRequest, 30),
			urls:        make([]tinyUrl, 0),
			rnd:         rand.New(rand.NewSource(time.Now().UnixNano())),
		}

		go processRequests(instance)

	}

	return instance
}

func (c *cache) Close() {
	close(c.done)
}

func processRequests(c *cache) {
	for {
		select {
		case t := <-c.writeBuffer:
			c.urls = append(c.urls, t)
		case r := <-c.readBuffer:
			size := len(c.urls)
			if size == 0 {
				//send request to the buffer again
				//because we do not have any url
				go func() {
					c.readBuffer <- r
				}()
			} else {
				rndIndex := c.rnd.Intn(size)
				r.response <- c.urls[rndIndex]
			}
		case <-c.done:
			return
		}
	}
}

func (c *cache) Write(tiny string, originalUrl string) {

	c.writeBuffer <- tinyUrl{
		tiny:        tiny,
		originalUrl: originalUrl,
	}

}

func (c *cache) ReadRandom() (string, string) {

	response := make(chan tinyUrl, 1)

	c.readBuffer <- readRequest{
		response,
	}

	res := <-response
	return res.tiny, res.originalUrl
}
