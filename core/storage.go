package core

// Storage represents the underlying storage for storing urls
type Storage interface {
	Put(string, string) error      //Put short url + original url
	Get(string) (string, error)    //Get original url by short url
	Visit(string, VisitInfo) error //Store visit information for short url
	Count(string) int              //Return visit count for specific short url
}
