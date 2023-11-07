package controller

import (
	"time"
)

func Today() int {
	now := time.Now().Format("2006-01-02")
	numOfToday := len(RegexFind(now))
	return numOfToday
}

func Total() int {
	total := len(RegexFind(""))
	return total
}

func Graph() map[string]interface{} {
	var date []string
	var count []int

	format := "2006-01-02"
	last14day, _ := time.Parse(format, time.Now().AddDate(0, 0, -14).Format("2006-01-02"))
	today, _ := time.Parse(format, time.Now().Format("2006-01-02"))

	for d := last14day; d.Before(today.AddDate(0, 0, 1)); d = d.AddDate(0, 0, 1) {
        date = append(date, d.Format(format))
    }

	for _, dd := range date {
		count = append(count, len(RegexFind(dd)))
	}

	return map[string]interface{} {"date": date, "count": count}
}