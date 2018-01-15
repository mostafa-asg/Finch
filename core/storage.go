package core

// Storage represents the underlying storage for storing urls
type Storage interface {
	Put(string, string) (bool, error)
	Get(string) (string, error)
}
