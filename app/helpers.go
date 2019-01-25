// Package app stores the base application
package app

import (
    "github.com/gorilla/mux"
    "github.com/huangsam/keyauth/apis"
)

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
    apis.ApiEndpoints[t] = m
    return nil
}
