// Package endpoints stores API endpoints
package endpoints

import (
    "encoding/json"
    "fmt"
    "github.com/gorilla/handlers"
    "github.com/gorilla/mux"
    "net/http"
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

// GetRouter serves application
func GetRouter() http.Handler {
    r := mux.NewRouter()
    apir := r.PathPrefix("/api/").Subrouter()
    apir.Handle("/apikey/", ApiKeyCoarse).Methods("GET", "POST", "OPTIONS")
    apir.Handle("/apikey/{id}/", ApiKeyGranular).Methods("GET", "DELETE", "OPTIONS")
    apir.HandleFunc("/apikey/{id}/archive/", ArchiveApiKey).Methods("PATCH", "OPTIONS")
    apir.HandleFunc("/apikey/authenticate/", AuthenticateApiKey).Methods("POST", "OPTIONS")
    miscr := r.PathPrefix("/").Subrouter()
    miscr.HandleFunc("/health/", HealthCheck).Methods("GET")
    miscr.HandleFunc("/", GetEndpoints).Methods("GET")
    apir.Walk(registerEndpoints)
    miscr.Walk(registerEndpoints)
    headersOK := handlers.AllowedHeaders(allowedHeaders)
    originsOK := handlers.AllowedOrigins(allowedOrigins)
    methodsOK := handlers.AllowedMethods(allowedMethods)
    return handlers.CORS(headersOK, originsOK, methodsOK)(r)
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
    if r.Method == "OPTIONS" {
        return
    }
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
    if r.Method == "OPTIONS" {
        return
    }
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
