# keyauth

Golang server that manages API keys for authenticated users

## Getting started

### Package installation

First download the dependencies:

    export GOPATH="..."
    go get github.com/gorilla/handlers
    go get github.com/gorilla/mux

### Server startup

Then run the main package:

    go run main.go

Access the server from port 3000 using your client of choice.

### Server tests

The server was tested at the following key endpoints:

- `/`
- `/health/`
- `/api/apikey/`

More endpoints can be tested in the future. To run the suite:

    go test
