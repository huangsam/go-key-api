// Package main initializes a HTTP server for managing API keys. One
// enhancement for this solution is to add CORS. This would simply mean
// adding the CORS handler to the router and OPTION method to all endpoints.
package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"strconv"
	"time"
)

// ApiKey is uniquely identified by content
type ApiKey struct {
	Id           int       `json:"id"`
	UserId       int       `json:"user_id"`
	Content      string    `json:"api_key"`
	TimeCreated  time.Time `json:"time_created"`
	TimeLastUsed time.Time `json:"time_last_used"`
	TimeArchived time.Time `json:"time_archived"`
}

// ServerStatus has message and failures
type ServerStatus struct {
	Message  string   `json:"message"`
	Failures []string `json:"failures"`
}

// Local state
var apiEndpoints map[string][]string = make(map[string][]string)
var apiKeyNotFound string = "404 api key not found"
var httpRouter http.Handler

// Dummy keys
var apiKeys []ApiKey = []ApiKey{
	ApiKey{1, 17, "0001-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{2, 17, "0002-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{3, 18, "0003-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{4, 20, "0004-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{5, 20, "0005-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{6, 17, "0006-xx-yy", time.Time{}, time.Time{}, time.Time{}},
	ApiKey{7, 19, "0007-xx-yy", time.Time{}, time.Time{}, time.Time{}},
}

// Dummy sequence number
var apiKeySequenceNumber = 5

// registerEndpoints registers paths and associated methods
func registerEndpoints(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	t, err := route.GetPathTemplate()
	if err != nil {
		return err
	}
	m, err := route.GetMethods()
	if err != nil {
		return err
	}
	apiEndpoints[t] = m
	return nil
}

// findApiKey finds match by array index
func findApiKey(searchId int) (int, bool) {
	for index, item := range apiKeys {
		if item.Id == searchId {
			return index, true
		}
	}
	return -1, false
}

// findApiKeyByContent finds match by content
func findApiKeyByContent(content string) (int, bool) {
	for index, item := range apiKeys {
		if item.Content == content {
			return index, true
		}
	}
	return -1, false
}

// findApiKeys finds matches by query parameters
func findApiKeys(queryParams map[string][]string) []ApiKey {
	userIdParam, ok := queryParams["user_id"]
	if !ok {
		return apiKeys
	}
	userId, err := strconv.Atoi(userIdParam[0])
	if err != nil {
		return apiKeys
	}
	resultKeys := []ApiKey{}
	for index, item := range apiKeys {
		if item.UserId == userId {
			resultKeys = append(resultKeys, apiKeys[index])
		}
	}
	return resultKeys
}

func main() {
	httpRouter = GetRouter()
	s := &http.Server{
		Addr:         ":3000",
		Handler:      handlers.LoggingHandler(os.Stdout, httpRouter),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	s.ListenAndServe()
}

// GetRouter serves application
func GetRouter() http.Handler {
	r := mux.NewRouter()
	apir := r.PathPrefix("/api/").Subrouter()
	apir.Handle("/apikey/", ApiKeyCoarse).Methods("GET", "POST")
	apir.Handle("/apikey/{id}/", ApiKeyGranular).Methods("GET", "DELETE")
	apir.HandleFunc("/apikey/{id}/archive/", ArchiveApiKey).Methods("PATCH")
	apir.HandleFunc("/apikey/authenticate/", AuthenticateApiKey).Methods("POST")
	miscr := r.PathPrefix("/").Subrouter()
	miscr.HandleFunc("/health/", HealthCheck).Methods("GET")
	miscr.HandleFunc("/", GetEndpoints).Methods("GET")
	apir.Walk(registerEndpoints)
	miscr.Walk(registerEndpoints)
	return r
}

// ApiKeyCoarse handles key retrieval and creation
var ApiKeyCoarse = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetApiKeys(w, r)
	case "POST":
		CreateApiKey(w, r)
	}
})

// GetApiKeys gets all API keys
func GetApiKeys(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	apiKeys := findApiKeys(queryParams)
	json.NewEncoder(w).Encode(apiKeys)
}

// CreateApiKey creates a single API key
func CreateApiKey(w http.ResponseWriter, r *http.Request) {
	var apiKey ApiKey
	err := json.NewDecoder(r.Body).Decode(&apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	apiKey.Id = apiKeySequenceNumber
	apiKey.Content = fmt.Sprintf("%04d-xx-yy", apiKeySequenceNumber)
	apiKey.TimeCreated = time.Now()
	apiKeys = append(apiKeys, apiKey)
	apiKeySequenceNumber += 1
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)
}

// ApiKeyGranular handles key retrieval and deletion
var ApiKeyGranular = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetApiKey(w, r)
	case "DELETE":
		DeleteApiKey(w, r)
	}
})

// GetApiKey gets a single API key
func GetApiKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index, ok := findApiKey(searchId)
	if !ok {
		http.Error(w, apiKeyNotFound, http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(apiKeys[index])
}

// DeleteApiKey deletes a single API key
func DeleteApiKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index, ok := findApiKey(searchId)
	if !ok {
		http.Error(w, apiKeyNotFound, http.StatusNotFound)
		return
	}
	apiKeys = append(apiKeys[:index], apiKeys[index+1:]...)
	w.WriteHeader(http.StatusNoContent)
}

// ArchiveApiKey archives a single API key
func ArchiveApiKey(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	searchId, err := strconv.Atoi(params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index, ok := findApiKey(searchId)
	if !ok {
		http.Error(w, apiKeyNotFound, http.StatusNotFound)
		return
	}
	apiKeys[index].TimeArchived = time.Now()
	w.WriteHeader(http.StatusNoContent)
}

// AuthenticateApiKey confirms the existence of a single API key
func AuthenticateApiKey(w http.ResponseWriter, r *http.Request) {
	var apiKey ApiKey
	err := json.NewDecoder(r.Body).Decode(&apiKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	index, ok := findApiKeyByContent(apiKey.Content)
	if !ok {
		http.Error(w, apiKeyNotFound, http.StatusNotFound)
		return
	}
	apiKeys[index].TimeLastUsed = time.Now()
	w.WriteHeader(http.StatusNoContent)
}

// HealthCheck checks server status
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := ServerStatus{"OK", []string{}}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// GetEndpoints gets all endpoints
func GetEndpoints(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(apiEndpoints)
}
