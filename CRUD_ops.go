package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CRUD operations

func getAllRates(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(sampleRates)
}

//reject out of hand if overnight, over month, or over year
func getRate(w http.ResponseWriter, r *http.Request) {
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
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(item.Price)
					return
				}

			}
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8") //for error message
		w.WriteHeader(http.StatusNotFound)                          //code 404
		w.Write([]byte("Rate not available for requested times."))
		return
	} else {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Rate not available for requested times."))
		return
	}
}

func addRate(w http.ResponseWriter, r *http.Request) {
	var rate Rate
	_ = json.NewDecoder(r.Body).Decode(&rate)
	//because it is being decoded to Rate type, string will revert to empty if invalid type and price will be set to 0.
	//Presumably there is no free parking.
	if rate.Days == "" || rate.Times == "" || rate.Price == 0 {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusBadRequest) //code 400
		w.Write([]byte("Rate entered was potentially malformed. Please check data types."))
		return
	}
	sampleRates = append(sampleRates, rate)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rate)
	return
}

func adjustRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, rate := range sampleRates {
		if rate.Days == params["days"] && rate.Times == params["hours"] {
			sampleRates = append(sampleRates[:index], sampleRates[index+1:]...)
			var rate Rate
			_ = json.NewDecoder(r.Body).Decode(&rate)
			//same logic as in addRate
			if rate.Days == "" || rate.Times == "" || rate.Price == 0 {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusBadRequest) //code 400
				w.Write([]byte("Rate entered was potentially malformed. Please check data types."))
				return
			}
			sampleRates = append(sampleRates, rate)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rate)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest) //code 400 - could not find to remove
}

func removeRate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for index, rate := range sampleRates {
		if rate.Days == params["days"] && rate.Times == params["hours"] {
			sampleRates = append(sampleRates[:index], sampleRates[index+1:]...)
			json.NewEncoder(w).Encode(sampleRates)
			return
		}
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusNotFound) //code 404 -- could not find to remove
	w.Write([]byte("Unable to locate specified rate for deletion."))
}
