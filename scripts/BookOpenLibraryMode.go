package scripts

type BookOpenLib struct {
	Title    string   `json:"title"`
	Subjects []string `json:"subjects"`
	Key      string   `json:"key"`
	Authors  []struct {
		Type   string `json:"type"`
		Author struct {
			Key string `json:"key"`
		} `json:"author"`
	} `json:"authors"`
}

type Book struct {
	Title    string   `json:"title"`
	Subjects []string `json:"subjects"`
	Authors  []string `json:"author"`
	Key      string   `json:"key"`
}

type Author struct {
	Name string `json:"name"`
}
