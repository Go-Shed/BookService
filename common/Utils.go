package common

import (
	"fmt"
	"time"
)

func GetFormattedDate(t string) string {

	thisTime, err := time.Parse("2006-01-02T15:04:05Z07:00", t)

	if err != nil {
		return fmt.Sprint(t)
	}

	duration := time.Now().Local().Sub(thisTime)

	if duration.Hours() < 1 {
		return fmt.Sprintf("%.0f m", duration.Minutes())
	} else if duration.Hours() < 24 {
		return fmt.Sprintf("%.0f h", duration.Hours())
	} else {
		return thisTime.Format("02/01/2006")
	}
}
