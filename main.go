package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Rate model based on sample data
type Rate struct {
	Days  string `json:"days"`
	Times string `json:"times"`
	Price int    `json:"price"`
}

var sampleRates []Rate

func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/rates/", getAllRates).Methods("GET")
	router.HandleFunc("/rates/{startTime}/{endTime}", getRate).Methods("GET")
	router.HandleFunc("/rates/", addRate).Methods("POST")
	router.HandleFunc("/rates/{days}/{hours}", adjustRate).Methods("PUT")
	router.HandleFunc("/rates/{days}/{hours}", removeRate).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func main() {

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

	startServer()
}
