package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetApiKeys(t *testing.T) {
	endpoint := "/api/apikey/"
	req, _ := http.NewRequest("GET", endpoint, nil)
	resp := httptest.NewRecorder()
	httpRouter := GetRouter()
	httpRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("%v status: %v", endpoint, resp.Code)
	}
	var apiKeys []ApiKey
	err := json.NewDecoder(resp.Body).Decode(&apiKeys)
	if err != nil {
		t.Fatalf("%v error: %v", endpoint, err)
	} else if len(apiKeys) != 7 {
		t.Fatalf("%v error: %v keys", endpoint, len(apiKeys))
	}
	for _, item := range apiKeys {
		currentContent := fmt.Sprintf("%04d-xx-yy", item.Id)
		if item.UserId < 17 || item.UserId > 20 {
			t.Fatalf("%v error: invalid user %v", endpoint, item.UserId)
		} else if len(item.Content) != 10 {
			t.Fatalf("%v error: invalid content %v", endpoint, item.Content)
		} else if currentContent != item.Content {
			t.Fatalf("%v error: %v != %v", endpoint, currentContent, item.Content)
		}
	}
}

func TestHealthCheck(t *testing.T) {
	endpoint := "/health/"
	req, _ := http.NewRequest("GET", endpoint, nil)
	resp := httptest.NewRecorder()
	httpRouter := GetRouter()
	httpRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("%v - %v", endpoint, resp.Code)
	}
}

func TestGetEndpoints(t *testing.T) {
	endpoint := "/"
	req, _ := http.NewRequest("GET", endpoint, nil)
	resp := httptest.NewRecorder()
	httpRouter := GetRouter()
	httpRouter.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("%v - %v", endpoint, resp.Code)
	}
}
