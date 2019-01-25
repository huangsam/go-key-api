# keyauth

Golang server that manages API keys for authenticated users

## Getting started

### Package installation

First download the dependencies:

    export GOPATH="..."
    go get github.com/gorilla/handlers
    go get github.com/gorilla/mux
    go install

### Server startup

Then run the main package:

    keyauth

Access the server from port 3000 using your client of choice.

## Developers

### Test suite

The server was tested at the following key endpoints:

- `/`
- `/health/`
- `/api/apikey/`

More endpoints can be tested in the future. To run the suite:

    go test

## Implementation Notes

Some things to consider:

- Closely resembles the folder structure [here](https://github.com/qiangxue/golang-restful-starter-kit)
- Currently does not use database connections
- Uses a fake data structure to act as a database
- CORS is properly enabled for the API
