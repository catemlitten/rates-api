package main

import (
	"fmt"
	"log"
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
	layout := "2006-01-02T15:04:05Z07:00" //ISO format
	t, err := time.Parse(layout, timeStamp)

	if err != nil {
		log.Fatal(err)
	}

	t.Format("Mon Jan _2 15:04:05 2006")
	day := t.Weekday().String()
	hour := strconv.Itoa(t.Hour())
	minute := strconv.Itoa(t.Minute())
	times := make([]string, 3)
	times[0] = day
	times[1] = hour
	times[2] = minute
	return times
}

func parseRateTimes(rates string) []int {
	rateSlice := strings.Split(rates, "-")
	//in 0000-2300 format
	rateTimes := make([]int, 2)
	firstHr, err := strconv.Atoi(rateSlice[0])
	if err != nil {
		log.Fatal(err)
	}
	secondHr, secErr := strconv.Atoi(rateSlice[1])
	if secErr != nil {
		log.Fatal(secErr)
	}
	rateTimes[0] = firstHr
	rateTimes[1] = secondHr
	return rateTimes
}

func isOverlappingOrInvalid(start string, end string) bool {
	layout := "2006-01-02T15:04:05Z07:00" //ISO format
	tStart, sErr := time.Parse(layout, start)
	if sErr != nil {
		log.Fatal(sErr)
	}
	tEnd, eErr := time.Parse(layout, end)
	if eErr != nil {
		log.Fatal(eErr)
	}
	if tStart.Year() != tEnd.Year() || tStart.Day() != tEnd.Day() || tStart.Month() != tEnd.Month() {
		return true
	}
	//if identical they have no actually requested valid time frame
	if tStart == tEnd {
		return true
	}
	return false
}

//Go does not support method overloading nor default params so string was used
func compareHours(rateTime int, requestedTimes []string, comparator string) bool {
	timeCombine := strings.Join(requestedTimes, "") // make [9,15] into '915'
	//1000, 900, 1015, 915, 000
	if comparator == "start" {
		requestStart, startErr := strconv.Atoi(timeCombine)

		if startErr != nil {
			log.Fatal(startErr)
		}
		fmt.Println(requestStart, rateTime)
		if requestStart > rateTime {
			return true
		}
		return false
	} else {
		requestEnd, endErr := strconv.Atoi(requestedTimes[1])
		if endErr != nil {
			log.Fatal(endErr)
		}
		if requestEnd < rateTime {
			return true
		}
		return false
	}
}
