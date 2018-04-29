// Package app stores the base application
package app

import (
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "github.com/huangsam/keyauth/endpoints"
    "net/http"
)

// GetRouter serves application
func GetRouter() http.Handler {
    r := mux.NewRouter()
    apir := r.PathPrefix("/api/").Subrouter()
    apir.Handle("/apikey/", endpoints.ApiKeyCoarse).Methods("GET", "POST", "OPTIONS")
    apir.Handle("/apikey/{id}/", endpoints.ApiKeyGranular).Methods("GET", "DELETE", "OPTIONS")
    apir.HandleFunc("/apikey/{id}/archive/", endpoints.ArchiveApiKey).Methods("PATCH", "OPTIONS")
    apir.HandleFunc("/apikey/authenticate/", endpoints.AuthenticateApiKey).Methods("POST", "OPTIONS")
    miscr := r.PathPrefix("/").Subrouter()
    miscr.HandleFunc("/health/", endpoints.HealthCheck).Methods("GET")
    miscr.HandleFunc("/", endpoints.GetEndpoints).Methods("GET")
    apir.Walk(registerEndpoints)
    miscr.Walk(registerEndpoints)
    headersOK := handlers.AllowedHeaders(allowedHeaders)
    originsOK := handlers.AllowedOrigins(allowedOrigins)
    methodsOK := handlers.AllowedMethods(allowedMethods)
    return handlers.CORS(headersOK, originsOK, methodsOK)(r)
}
