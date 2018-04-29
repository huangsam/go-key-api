// Package endpoints stores API endpoints
package endpoints

import (
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
