package common

import (
	"strings"
)

func GetDateOnlyFromISO8601(date string) string {
	dateTime := date
	if strings.Contains(date, "T") {
		dateTime = strings.Split(date, "T")[0]
	}

	return dateTime
}

func GetTimeOnlyFromISO8601(date string) string {
	dateTime := date
	if strings.Contains(date, "T") {
		dateTime = strings.Split(date, "T")[1]
	}

	return dateTime
}
