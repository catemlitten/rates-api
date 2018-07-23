package main

import (
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/*
Utility functions used by other aspects of the API. Abstracted out for clarity of code and reuse.
*/

//CompareRateDays should slice out appropriate portion of date string and determine if found in rate days
func CompareRateDays(days string, requestDay string) bool {
	//date parsed from user as 'Saturday'/full name
	//days stored as mon/tues/weds/thurs/fri/sat/sun
	requestDay = strings.ToLower(requestDay)
	if requestDay != "thursday" && requestDay != "tuesday" {
		requestDay = requestDay[0:3]
	} else if requestDay == "thursday" {
		requestDay = requestDay[0:5]
	} else {
		requestDay = requestDay[0:4]
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
	rateTimes := make([]int, 2)
	firstHr, err := strconv.Atoi(rateSlice[0])
	if err != nil {
		log.Fatal(err)
	}
	secondHr, err := strconv.Atoi(rateSlice[1])
	if err != nil {
		log.Fatal(err)
	}
	rateTimes[0] = firstHr
	rateTimes[1] = secondHr
	return rateTimes
}

/*
As noted in the specification the rates do not overlap and requests are to be within a timeframe.
As such any rates which can be determined quickly to be overlapping or invalid should be discarded.
*/
func isOverlappingOrInvalid(start string, end string) bool {
	layout := "2006-01-02T15:04:05Z07:00" //ISO format
	tStart, err := time.Parse(layout, start)
	if err != nil {
		//if there is an error parsing the date it is invalid
		return true
	}
	tEnd, err := time.Parse(layout, end)
	if err != nil {
		return true
	}
	if tStart.Year() != tEnd.Year() || tStart.Day() != tEnd.Day() || tStart.Month() != tEnd.Month() {
		return true
	}
	//if identical they have not actually requested valid time frame
	if tStart == tEnd {
		return true
	}
	return false
}

//Go does not support method overloading nor default params. Bool is true for 'start'.
func compareHours(rateTime int, requestedTimes []string, isStart bool) bool {
	requestHour := fixTime(requestedTimes)
	if isStart == true {
		if requestHour > rateTime {
			return true
		}
		return false
	}
	if requestHour < rateTime {
		return true
	}
	return false
}

func fixTime(requestedTimes []string) int {
	if requestedTimes[1] == "0" {
		requestedTimes[1] = "00" //so 900 is compared to 915 and not 90 v 915
	}
	if len(requestedTimes[1]) == 1 {
		requestedTimes[1] = "0" + requestedTimes[1]
	}
	timeCombine, err := strconv.Atoi(strings.Join(requestedTimes, "")) // make [9,15] into '915'
	if err != nil {
		log.Fatal(err)
	}
	return timeCombine
}

//Remove repetitive code within body of CRUD methods
func setNonSuccessHeader(w http.ResponseWriter, status int, details string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	w.Write([]byte(details))
}

//basic check for data quality in put/post
func rateTimeCheck(days string, times string) bool {
	containsDays := regexp.MustCompile(`mon|tues|wed|thurs|fri|sat|sun`)
	containsNumbers := regexp.MustCompile(`\d`)
	containsOneHyphen := regexp.MustCompile(`-{1}`)
	if !containsDays.MatchString(days) || containsNumbers.MatchString(days) || containsOneHyphen.MatchString(days) {
		return false
	}
	if containsDays.MatchString(times) || !containsNumbers.MatchString(times) || !containsOneHyphen.MatchString(times) {
		return false
	}
	return true
}

func setCORS(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "null") //demo purposes only. page will be static.
}
