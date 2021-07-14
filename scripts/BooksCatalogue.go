package scripts

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func CleanData(in string) error {

	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	csvr := csv.NewReader(f)
	csvr.Comma = '\t'
	csvr.FieldsPerRecord = -1
	csvr.LazyQuotes = true

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		var item BookOpenLib
		json.Unmarshal([]byte(row[4]), &item)

		book := Book{
			Title:    item.Title,
			Key:      item.Key,
			Subjects: item.Subjects,
		}

		var authors []string
		for _, author := range item.Authors {
			authorName, _ := getAuthor(author.Author.Key)
			authors = append(authors, authorName)
		}
		book.Authors = authors

		result, _ := json.Marshal(book)

		fmt.Println(string(result))
	}
}

func getAuthor(key string) (string, error) {

	URL := "https://openlibrary.org/" + key + ".json"

	request, err := http.NewRequest("GET", URL, nil)

	if err != nil {
		return "none", err
	}

	client := &http.Client{Timeout: time.Second * 10}
	response, err := client.Do(request)

	if err != nil {
		return "none", err
	}
	defer response.Body.Close()
	result, _ := ioutil.ReadAll(response.Body)

	var author Author
	json.Unmarshal([]byte(result), &author)

	return author.Name, nil
}
