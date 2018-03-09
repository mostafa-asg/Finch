package core

//Stats contains aggregate statistics about url
type Stats struct {
	Timeline  map[string]int
	Referrals map[string]int
	Browsers  map[string]int
	Countries map[string]int
	OS        map[string]int
}
