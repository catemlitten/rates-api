package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func CompareRateDays(days string, requestDay string) bool {
	//date parsed from user as 'Saturday'/full name
	//days stored as mon/tues/weds/thurs/fri/sat/sun
	requestDay = strings.ToLower(requestDay)
	if requestDay != "thursday" {
		requestDay = requestDay[0:3]
	} else {
		requestDay = requestDay[0:5]
	}
	return strings.Contains(days, requestDay)
}

func timeStampRead(timeStamp string) []string {
	layout := "2006-01-02T15:04:05Z07:00"
	t, err := time.Parse(layout, timeStamp)

	if err != nil {
		fmt.Println(err)
	}

	t.Format("Mon Jan _2 15:04:05 2006")
	day := t.Weekday().String()
	hour := strconv.Itoa(t.Hour())
	minute := strconv.Itoa(t.Minute())
	sec := strconv.Itoa(t.Second())
	times := make([]string, 4)
	times[0] = day
	times[1] = hour
	times[2] = minute
	times[3] = sec
	return times
}

func parseRateTimes(rates string) []int {
	// rates = strconv.Atoi(rates)
	rateSlice := strings.Split(rates, "-")
	//in 0000-2300 format
	rateTimes := make([]int, 2)
	firstHr, err := strconv.Atoi(rateSlice[0])
	if err != nil {
		fmt.Println(err)
	}
	secondHr, err := strconv.Atoi(rateSlice[1])
	if err != nil {
		fmt.Println(err)
	}
	rateTimes[0] = firstHr
	rateTimes[1] = secondHr
	return rateTimes
}
