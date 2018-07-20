package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// general mux/http testing for golang based off tutorial: https://www.thepolyglotdeveloper.com/2017/02/unit-testing-golang-application-includes-http/
func Router() *mux.Router {
	go startServer()
	router := mux.NewRouter()
	router.HandleFunc("/rates/", getAllRates).Methods("GET")
	router.HandleFunc("/rates/{startTime}/{endTime}", getRate).Methods("GET")
	router.HandleFunc("/rates/", addRate).Methods("POST")
	router.HandleFunc("/rates/{days}/{hours}", adjustRate).Methods("PUT")
	router.HandleFunc("/rates/{days}/{hours}", removeRate).Methods("DELETE")
	return router
}

func TestGetAll(t *testing.T) {
	testTable := []struct {
		name   string
		status int
	}{
		{name: "Testing generic GET, no params", status: 200},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)
			request, _ := http.NewRequest("GET", "/rates/", nil)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
		})
	}
}

func TestGetSelectRate(t *testing.T) {
	testTable := []struct {
		name      string
		startTime string
		endTime   string
		rate      int
		status    int
		err       string
	}{
		{name: "Testing valid rate: Wednesday 1:15AM to 4am", startTime: "2015-07-01T01:15:00Z", endTime: "2015-07-01T04:00:00Z", rate: 1000, status: 200},
		{name: "Testing unavailable rate: Wednesday 1:15am to 5am (is not inclusive)", startTime: "2015-07-01T01:15:00Z", endTime: "2015-07-01T05:00:00Z", err: "Rate not available for requested times.", status: 404},
		{name: "Testing bad params: No end time", startTime: "2015-07-01T01:15:00Z", err: "404 page not found", status: 404},
		{name: "Testing bad params: Invalid dates (bad format)", startTime: "abcedefg", err: "404 page not found", status: 404},
		{name: "Testing bad params: Invalid dates (times outside range)", startTime: "2015-07-01T01:15:00Z", endTime: "2018-07-01T05:00:00Z", err: "Rate not available for requested times.", status: 404},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)
			request, _ := http.NewRequest("GET", "/rates/"+testCase.startTime+"/"+testCase.endTime, nil)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			body := 0
			json.NewDecoder(response.Body).Decode(&body)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
			if testCase.status == 200 {
				assert.Equal(t, testCase.rate, body, "Matching response is expected")
			}
		})
	}
}

func TestAddRate(t *testing.T) {
	testTable := []struct {
		name   string
		days   string
		times  string
		rate   int
		status int
	}{
		{name: "Testing valid rate: Thursday 4:00AM to 9:00am, rate of 2500", days: "thurs", times: "0400-0900", rate: 2500, status: 200},
		{name: "Testing valid rate: Monday, Thursday, Saturday 10:00AM to 9:00PM, rate of 370", days: "mon,thurs,sat", times: "1000-2100", rate: 3700, status: 200},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)
			rate := &Rate{
				Days:  testCase.days,
				Times: testCase.times,
				Price: testCase.rate,
			}
			jsonRateObj, _ := json.Marshal(rate)
			request, _ := http.NewRequest("POST", "/rates/", bytes.NewBuffer(jsonRateObj))
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
		})
	}
}

/* unclear how to fully replicate behavior of live app.
Days and times which are not strings become "" and rates which are not ints become
the number zero. Due to strict typing however it was not possible to feed incorrect data types into the rate without
manually converting first.
*/
func TestMalformedPostData(t *testing.T) {
	testTable := []struct {
		name   string
		days   string
		times  string
		rate   int
		status int
	}{
		{name: "Testing badly formatted data (days)", days: "", times: "0900-0700", rate: 100, status: 400},
		{name: "Testing badly formatted data (hours)", days: "mon,tues", times: "", rate: 100, status: 400},
		{name: "Testing badly formatted data (price)", days: "mon,tues", times: "0900-0700", rate: 0, status: 400},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)

			rate := &Rate{
				Days:  testCase.days,
				Times: testCase.times,
				Price: testCase.rate,
			}
			jsonRateObj, _ := json.Marshal(rate)
			request, _ := http.NewRequest("POST", "/rates/", bytes.NewBuffer(jsonRateObj))
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
		})
	}
}

func TestAdjustRate(t *testing.T) {
	testTable := []struct {
		name     string
		adjDays  string
		adjTimes string
		newDays  string
		newTimes string
		rate     int
		status   int
	}{
		{name: "Adjusting exisiting time: Wednesday 6am-6pm to Tu/Th 6am-Noon", adjDays: "wed", adjTimes: "0600-1800", newDays: "tues,thurs", newTimes: "0600-1200", rate: 2500, status: 200},
		{name: "Attempt to adjust non-existant rate Mon/Thur/Sat from 10am-9pm", adjDays: "mon,thur,sat", adjTimes: "1000-2100", newDays: "mon,thurs", newTimes: "1100-1400", rate: 3700, status: 400},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)
			rate := &Rate{
				Days:  testCase.newDays,
				Times: testCase.newTimes,
				Price: testCase.rate,
			}
			jsonRateObj, _ := json.Marshal(rate)
			request, _ := http.NewRequest("PUT", "/rates/"+testCase.adjDays+"/"+testCase.adjTimes, bytes.NewBuffer(jsonRateObj))
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
		})
	}
}

func TestRemoveRate(t *testing.T) {
	testTable := []struct {
		name   string
		days   string
		times  string
		rate   int
		status int
	}{
		{name: "Remove exisiting time: Wednesday 6am-6pm to Sun/Tu 1am-7am", days: "sun,tues", times: "0100-0700", status: 200},
		{name: "Attempt to remove non-existant rate Mon/Thur/Sat from 10am-9pm", days: "mon,thur,sat", times: "1000-2100", status: 404},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			fmt.Println(testCase.name)
			request, _ := http.NewRequest("DELETE", "/rates/"+testCase.days+"/"+testCase.times, nil)
			response := httptest.NewRecorder()
			Router().ServeHTTP(response, request)
			assert.Equal(t, testCase.status, response.Code, "Matching response is expected")
		})
	}
}
