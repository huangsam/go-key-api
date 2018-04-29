// Package main runs a key authentication server
package main

import (
	"github.com/gorilla/handlers"
	"github.com/huangsam/keyauth/endpoints"
	"net/http"
	"os"
	"time"
)

// Local state
var httpRouter http.Handler

func main() {
	httpRouter = endpoints.GetRouter()
	s := &http.Server{
		Addr:         ":3000",
		Handler:      handlers.LoggingHandler(os.Stdout, httpRouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.ListenAndServe()
}
