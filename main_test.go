package main

import (
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
