package core

//VisitInfo encapsulates the information of each visit
type VisitInfo struct {
	Year     int
	Month    int
	Day      int
	Hour     int
	Minute   int
	Referrer string
	Browser  string
	Country  string
	OS       string
}
