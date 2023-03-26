# Accepted Home Project

---

## Description

Service that provides a REST API offering CRUD operations will be enable clients to receive, store data (json files) and
generate csv files

---

## Run

`cd script && make start-app`

* This command will start the app with `localhost` address and `:8080` port (specified in build/Dockerfile.devand .env) 

---

## Makefile Commands

| Command                         | Usage                                                                  |
|---------------------------------|------------------------------------------------------------------------|
| start-app                       | `Start app`                                                            |
| kill-app                        | `Stop app`                                                             |
| rebuild-app                     | `Rebuild app`                                                          |
| tests-all                       | `Run both unit and integration tests`                                  |
| tests-benchmark                 | `Run benchmark tests`                                                  |
| tests-unit                      | `Run unit tests `                                                      |
| tests-file FILE={filePath}      | `Run specific file test`                                               |
| generate-mock FILE={filePath}   | `Generate mock for a specific file`                                    |
| tests-package PACKAGE={package} | `Run specific package test`                                            |
| tests-all-with-coverage         | `Run both unit and integration tests via docker with coverage details` |

* All these are executed through docker containers
* In order to execute makefile commands type **make** plus a command from the table above

  make {command}

---

## Endpoints

1. ### Create Solar Panel Data

POST /solar-panel-data

#### Request

```json
{
  "solar": {
    "38d503e5-dc1c-4549-8172-09d9c29070f7": [
      [
        "20211231T221500Z",
        "0.0"
      ]
    ]
  },
  "wind": null
}
```

#### Response

##### Success

Status Code *201 Created*

```json
{
  "id": "0e96297f-ad56-426f-864e-5ac3aca5c3e7",
  "dataSubmitted": {
    "solar": {
      "38d503e5-dc1c-4549-8172-09d9c29070f7": [
        [
          "20211231T221500Z",
          "0.0"
        ]
      ]
    },
    "wind": null
  }
}
```

##### Failure

Status Code *400 Bad Request* for malformed json or missing solar data  
Status Code *500 Interval Server Error*

2. ### Read Solar Panel Data

GET /solar-panel-data/{uuid}

#### Request

#### Response

##### Success

Status Code *200 OK*

```csv
Events
0.0
```

##### Failure

Status Code *404 Not Found Request* for not existing uuid  
Status Code *500 Interval Server Error*

3. ### Update Solar Panel Data

PUT /solar-panel-data/{uuid}

#### Request

```json
{
  "solar": {
    "38d503e5-dc1c-4549-8172-09d9c29070f7": [
      [
        "20211231T221500Z",
        "0.0"
      ]
    ]
  },
  "wind": null
}
```

#### Response

##### Success

Status Code *200 OK*

##### Failure

Status Code *400 Bad Request* for malformed json or missing solar data  
Status Code *404 Not Found Request* for not existing uuid  
Status Code *500 Interval Server Error*

4. ### Delete Solar Panel Data

DELETE /solar-panel-data/{uuid}

#### Request

#### Response

##### Success

Status Code *200 OK*

##### Failure

Status Code *404 Not Found Request* for not existing uuid  
Status Code *500 Interval Server Error*

---

## Notes

1. .env is pushed to Git only for the Assessment purpose. Config and .env files should never be tracked.
2. There are three Dockerfile files.
    1. Dockerfile is the normal, production one
    2. Dockerfile.dev is for setting up a remote debugger Delve
    3. Dockerfile.utilities is for building a docker for "utilities" like running tests,  
       linting etc
3. As I was implementing this, I noticed that package `github.com/gorilla/mux` is no longer  
   maintained. So normally, I would search for an alternative router. Due to lack of time
   I did not.
4. I had a choice to make on whether I should return the data on GET if one or more  
   parameters and events are malformed (for example if an event misses a value). This is  
   a choice I would make with the Product/Business team depending on if this missing value  
   corrupts the result or not if missing. For this exercise, I considered that if an event  
   or more are missing, then I can't guarantee the validity of the data so I'm returning an error.
5. I had a doubt about whether the Update operation should be an Upsert operation (which means  
   to create the element if not exists). PUT http verb in RESTful design supports both, so it's
   just a matter of choice. I chose to return an error if the data does not exist because    
   that's what I understood from the project specifications
6. Solar data validation is very simple, only check if empty. Should have done more, but had no time 
7. Graceful Shutdown does not work, did not have time to fix it