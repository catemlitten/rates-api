package main

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
)

//TestGetAll determines if generic GET works
func TestGetAll(t *testing.T) {
	go startServer()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//https://godoc.org/github.com/jarcoal/httpmock#RegisterResponder
	httpmock.RegisterResponder("GET", "http://localhost:8080/rates/",
		httpmock.NewStringResponder(200, `[
			{
			Days:  "mon,tues,thurs",
			Times: "0900-2100",
			Price: 1500
		},
			]`))

	r, _ := http.NewRequest("GET", "http://localhost:8080/rates/", nil)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGetRate(t *testing.T) {
	go startServer()
	client := &http.Client{
		Timeout: 1 * time.Second,
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	//https://godoc.org/github.com/jarcoal/httpmock#RegisterResponder
	httpmock.RegisterResponder("GET", "http://localhost:8080/rates/",
		httpmock.NewStringResponder(200, `[
			{
			Days:  "mon,tues,thurs",
			Times: "0900-2100",
			Price: 1500
		},
		{
			Days:  "fri,sat,sun",
			Times: "0900-2100",
			Price: 2000
		},
		{
			Days:  "wed",
			Times: "0600-1800",
			Price: 1750
		},
		{
			Days:  "mon,wed,sat",
			Times: "0100-0500",
			Price: 1000
		},
		{
			Days:  "sun,tues",
			Times: "0100-0700",
			Price: 925
		}
			]`))

	r, _ := http.NewRequest("GET", "http://localhost:8080/rates/2015-07-04T09:10:00Z/2015-07-04T20:59:00Z", nil)

	resp, err := client.Do(r)
	if err != nil {
		panic(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
