package scripts

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
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

		result, _ := json.Marshal(item)

		fmt.Println(string(result))
	}
}
