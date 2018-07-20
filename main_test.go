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
