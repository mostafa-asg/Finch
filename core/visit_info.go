package core

import (
	"time"
)

//VisitInfo encapsulates the information of each visit
type VisitInfo struct {
	Time     time.Time
	Referrer string
	Browser  string
	Country  string
	OS       string
}
