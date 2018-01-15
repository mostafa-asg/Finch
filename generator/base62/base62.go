package base62

import (
	"math/rand"
	"time"
)

type base62Gen struct {
	length int
}

type conBase62Gen struct {
}

var generateRequestChan = make(chan string, 5)
var generateResponseChan = make(chan string, 5)
var doneChan = make(chan bool)

var chars = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var random = rand.New(rand.NewSource(time.Now().UnixNano()))

const newGenerateRequest = "R"

// NewConcurrentBase62 generates random base62 string
func NewConcurrent() *conBase62Gen {

	for i := 2; i <= 6; i++ {
		go generateIDAgent(i)
	}

	return &conBase62Gen{}
}

func generateIDAgent(length int) string {

	g := New(length)

	for {

		select {
		case <-generateRequestChan:
			generateResponseChan <- g.GenerateID()
		case <-doneChan:
			break
		}

	}

}

// New generates random base62 string
func New(length int) *base62Gen {
	return &base62Gen{
		length,
	}
}

func (g *base62Gen) GenerateID() string {

	result := make([]byte, g.length)

	for i := 0; i < g.length; i++ {
		result[i] = chars[random.Intn(len(chars))]
	}

	return string(result)
}

func (g *conBase62Gen) GenerateID() string {
	generateRequestChan <- newGenerateRequest
	return <-generateResponseChan
}
