# Rate API

**Rate API** is a REST API written in Go. Given a set of dates written in ISO format it will retrieve, if available, the pricing for the given interval. There is an assumption that rates will never overlap.


Table of Contents
-----------------

[**Specifications**](#specs)

[**End Points**](#endpoints)

[**Build & Use Instructions**](#build)

<br />

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
The user must be able to curl against the API and recieve a price. For example, `curl http://localhost:8080/rates/2018-07-19T09:10:00Z/2018-07-19T20:59:00Z`checks for a rate on a Thursday between the hours of 9:10am and 8:59pm and would recieve back `1500`.

<a name="endpoints"></a>
End Points
------------

### /rates/

A user can GET against `/rates/` without parameters and will recieve back all rates.

A user can POST against `/rates/` to insert a new rate. As the data is static this will not persist after shutdown.

### /rates/{start}/{end}

A user can GET with a start and end time to recieve a specific price. A start date without an end date will return a 404.

### /rates/{days}/{hours}

A user can PUT or DELETE to this endpoint and either update a a specified rate or remove it. As the data is static this will not persist after shutdown.

<a name="build"></a>
Build and Use
--------------------------

`cd` into `rates-api` and run `go build` followed by `./rates-api.exe`. You can either use curl or Postman to test the API as it has dummy data loaded in.