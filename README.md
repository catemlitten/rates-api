# Rate API

**Rate API** is a REST API written in Go. Given a set of dates written in ISO format it will retrieve, if available, the pricing for the given interval. There is an assumption that rates will never overlap.


Table of Contents
-----------------

[**Specifications**](#specs)

[**Build & Use Instructions**](#build)

[**Testing**](#testing)


<a name="specs"></a>
Specifications
--------------------------

Data will be in JSON with the following structure:
```
       {   
            "days": "mon,tues,thurs",
            "times": "0900-2100",
            "price": 1500
        }
    
```
The user must be able to curl against the API and recieve a price. For example, `curl -i http://localhost:8080/rates/2018-07-19T09:10:00Z/2018-07-19T20:59:00Z`checks for a rate on a Thursday between the hours of 9:10am and 8:59pm and would recieve back `1500` with a status code of `200`. `curl -i http://localhost:8080/rates/2018-07-19T01:09:00Z/2018-07-19T21:00:00Z` would return `Rate not available for requested times.` and a staus code of `404`.

The full API contract is located in */Contract* and was written using <a href="https://apiblueprint.org/">API Blueprint</a>. <a href="https://www.npmjs.com/package/aglio">aglio</a> was used to create the end-user friendly *contract.html*

The desired implementations of *POST*, *PUT*, and *DELETE* operations were not explicitly detailed in the specifications document, so general use assumptions were made.

<a name="build"></a>
Build and Use
--------------------------

`cd` into `rates-api` and run `go get ./..` (this will install all dependencies) followed by `go build`. Once the program is built run `./rates-api.exe`. You can either use curl or Postman to test the API as it has dummy data loaded in.

<a name="testing"></a>
Testing
--------------------------

Testing was done using the <a href="https://golang.org/pkg/net/http/httptest/">httptest</a> package as well as the <a href="https://github.com/stretchr/testify">testify</a> toolkit. As this was a first attempt doing API testing in Go, it is possible that tests were not as robust as hoped for. Currently it starts the Mux server and relies on the same mocked data that a live usage would have access to. After building the program as noted above, type `go test` into the terminal and all tests will run.