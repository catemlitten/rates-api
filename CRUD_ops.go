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

func getRate(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	//reject out of hand if overnight, over month, or over year
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
		setNonSuccessHeader(w, http.StatusNotFound, "Rate not available for requested times.")
		return
	}
	setNonSuccessHeader(w, http.StatusNotFound, "Rate not available for requested times.")
	return
}

func addRate(w http.ResponseWriter, r *http.Request) {
	var rate Rate
	_ = json.NewDecoder(r.Body).Decode(&rate)
	//because it is being decoded to Rate type, string will revert to empty if invalid type and price will be set to 0.
	//Presumably there is no free parking.
	if rate.Days == "" || rate.Times == "" || rate.Price == 0 {
		setNonSuccessHeader(w, http.StatusBadRequest, "Rate entered was potentially malformed. Please check data types.")
		return
	}
	if !rateTimeCheck(rate.Days, rate.Times) {
		setNonSuccessHeader(w, http.StatusBadRequest, "Rate entered was potentially malformed. Please check data types.")
		return
	}
	sampleRates = append(sampleRates, rate)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
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
			if rate.Days == "" || rate.Times == "" || rate.Price == 0 {
				setNonSuccessHeader(w, http.StatusBadRequest, "Rate entered was potentially malformed. Please check data types.")
				return
			}
			if !rateTimeCheck(rate.Days, rate.Times) {
				setNonSuccessHeader(w, http.StatusBadRequest, "Rate entered was potentially malformed. Please check data types.")
				return
			}
			sampleRates = append(sampleRates, rate)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(rate)
			return
		}
	}
	setNonSuccessHeader(w, http.StatusNotFound, "Could not locate rate to adjust.")
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
	setNonSuccessHeader(w, http.StatusNotFound, "Unable to locate specified rate for deletion.")
}
