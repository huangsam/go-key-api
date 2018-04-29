// Package endpoints stores API endpoints
package endpoints

import (
    "time"
)

// Local state
var apiEndpoints map[string][]string = make(map[string][]string)

// Constants
var apiKeyNotFound string = "404 api key not found"
var allowedHeaders []string = []string{"Authorization", "Content-Type", "X-Requested-With"}
var allowedOrigins []string = []string{"*"}
var allowedMethods []string = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}

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
var apiKeySequenceNumber = len(apiKeys) + 1
