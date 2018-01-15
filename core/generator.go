package core

// Generator generates random string
// This string later will be used as a row ID in database
type Generator interface {
	GenerateID() string
}
