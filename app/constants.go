// Package app stores the base application
package app

var allowedHeaders []string = []string{"Authorization", "Content-Type", "X-Requested-With"}
var allowedOrigins []string = []string{"*"}
var allowedMethods []string = []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}
