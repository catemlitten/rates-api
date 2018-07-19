package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Use structs as models
type Rate struct {
	Days  string `json:"days"`
	Times string `json:"times"`
	Price int    `json:"price"`
}

var sampleRates []Rate

//Verbs
func getAllRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sampleRates)
}

//reject out of hand if overnight, over month, or over year
func getRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	if !isOverlappingOrInvalid(params["startTime"], params["endTime"]) {
		requestStart := timeStampRead(params["startTime"])
		requestEnd := timeStampRead(params["endTime"])

		requestDay := requestStart[0]
		requestStartTime := requestStart[1:]
		requestEndTime := requestEnd[1:]

		for _, item := range sampleRates {
			if CompareRateDays(item.Days, requestDay) {
				//check to see if times requested available
				hours := parseRateTimes(item.Times)
				startHr := hours[0]
				endHr := hours[1]
				if compareHours(startHr, requestStartTime, "start") && compareHours(endHr, requestEndTime, "end") {
					json.NewEncoder(w).Encode(item)
					return
				}

			}
		}
		json.NewEncoder(w).Encode(&Rate{})
	} else {
		json.NewEncoder(w).Encode("No rates available for specified time frame")
		return
	}
}

func addRate(w http.ResponseWriter, r *http.Request) {

}

func adjustRate(w http.ResponseWriter, r *http.Request) {

}

func removeRate(w http.ResponseWriter, r *http.Request) {

}

/*
User should be able to curl against API with iso format dates and get back rates if available, notice of unavailable if do not exist
TODO:
Convert ISO dates to Mon-Sun dates, times. Return rate if avaialble
If start time === request start time or vv endtime === request end time, does not fully encapuslate and is not available

get rates by time frame
parse JSON for days => if day match determine if time match.
*/

func main() {
	router := mux.NewRouter()

	//Sample data
	sampleRates = append(sampleRates, Rate{
		Days:  "mon,tues,thurs",
		Times: "0900-2100",
		Price: 1500})
	sampleRates = append(sampleRates, Rate{
		Days:  "fri,sat,sun",
		Times: "0900-2100",
		Price: 2000})
	sampleRates = append(sampleRates, Rate{
		Days:  "wed",
		Times: "0600-1800",
		Price: 1750})
	sampleRates = append(sampleRates, Rate{
		Days:  "mon,wed,sat",
		Times: "0100-0500",
		Price: 1000})
	sampleRates = append(sampleRates, Rate{
		Days:  "sun,tues",
		Times: "0100-0700",
		Price: 925})

	router.HandleFunc("/rates/", getAllRates).Methods("GET")
	router.HandleFunc("/rates/{startTime}/{endTime}", getRate).Methods("GET")
	router.HandleFunc("/rates", addRate).Methods("POST")
	router.HandleFunc("/rates/{timeFrame}", adjustRate).Methods("PUT")
	router.HandleFunc("/rates/{timeFrame}", removeRate).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}
