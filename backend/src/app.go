/**
 * app.go
 * This file implements initial script of Slack Assistant Yui.
 *
 * auther u-mochi
 * license MIT
 */

// Package yui implements initial script of Slack Assistant Yui.
package yui

import (
	"net/http"
	"rest"
	"todoist"

	"appengine"
)

const (
	api = "api"
)

// init initializes this package
func init() {
	rest.RegieterMapping(api, "todoist", "configuration", todoist.ReadOperation, rest.RespondMethodNotAllowed, todoist.UpdateOperation, rest.RespondMethodNotAllowed)
	http.HandleFunc("/", handler)
}

// handler processes requests
func handler(writer http.ResponseWriter, request *http.Request) {
	rest.Process(writer, request, appengine.NewContext(request))
}
