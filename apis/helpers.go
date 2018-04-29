// Package apis stores API endpoints
package apis

import (
    "strconv"
)

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
