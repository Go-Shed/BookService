package scripts

type BookOpenLib struct {
	NumberOfPages int      `json:"number_of_pages"`
	Isbn10        []string `json:"isbn_10"`
	Series        []string `json:"series"`
	Key           string   `json:"key"`
	Authors       []struct {
		Key string `json:"key"`
	} `json:"authors"`
	Genres []string `json:"genres"`
	Title  string   `json:"title"`
}
