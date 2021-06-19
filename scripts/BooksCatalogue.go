package scripts

import (
	"encoding/csv"
	"io"
	"os"
)

func CleanData(in, out string) error {

	f, err := os.Open(in)
	if err != nil {
		return err
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			return err
		}

		item := UpdateWithdrawlRequest{WithdrawalRequestId: row[0]}
		response = append(response, item)
	}

}

func (service *WithdrawlService) csvTOList(filename string) ([]UpdateWithdrawlRequest, error) {

}
