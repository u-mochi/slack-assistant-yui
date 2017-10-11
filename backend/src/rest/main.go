/**
 * rest/main_test.go
 * This file implements REST framework of rest package
 *
 * auther u-mochi
 * license MIT
 */

// Package rest implements REST framework of Slack Assistant Yui.
package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"google.golang.org/appengine/log"
)

const (
	// MethodGet means GET request
	MethodGet = "GET"
	// MethodPost means POST request
	MethodPost = "POST"
	// MethodPut means PUT request
	MethodPut = "PUT"
	// MethodDelete means DELETE request
	MethodDelete = "DELETE"
	// MethodHead means HEAD request
	MethodHead = "HEAD"
	// MethodOption means OPTION request
	MethodOption = "OPTION"
	// URLPathSeparator means separator of path in URL
	URLPathSeparator = "/"
	// ErrorURLMappingUnMatch means URL mapping is not match.
	ErrorURLMappingUnMatch = "URL mapping is not match."
)

var (
	// urlMappings is slice of urlMapping structures
	urlMappings []urlMapping
)

// urlMapping means mapping of URL
type urlMapping struct {
	// Prefix of URL (ie. api of /api/todoist/configure/somekey)
	prefix string
	// Name of module (ie. todoist of /api/todoist/configure/somekey)
	moduleName string
	// Name of Model in URL (ie. configure of /api/todoist/configure/somekey)
	modelName string
	// Identifier of Model in URL (ie. configure of /api/todoist/configure/somekey)
	identifier string
	// Function on GET Method
	funcGet func(context.Context, string, http.ResponseWriter, *http.Request)
	// Function on POST Method
	funcPost func(context.Context, string, http.ResponseWriter, *http.Request)
	// Function on PUT Method
	funcPut func(context.Context, string, http.ResponseWriter, *http.Request)
	// Function on DELETE Method
	funcDelete func(context.Context, string, http.ResponseWriter, *http.Request)
}

// Initialize this package
func init() {
	urlMappings = make([]urlMapping, 0, 10)
}

// selectMapping selects registered urlMapping matches specified URL
func selectMapping(context context.Context, url *url.URL) (*urlMapping, string, error) {
	paths := strings.Split(strings.Trim(url.Path, URLPathSeparator), URLPathSeparator)
	if len(paths) > 2 {
		for _, mapping := range urlMappings {
			if mapping.prefix == paths[0] && mapping.moduleName == paths[1] && mapping.modelName == paths[2] {
				if len(paths) > 3 {
					return &mapping, paths[3], nil
				}
				return &mapping, "", nil
			}
		}
	}
	log.Infof(context, "mapping of %s not found in %#v", url.Path, urlMappings)
	return new(urlMapping), "", errors.New(ErrorURLMappingUnMatch)
}

// HandleRequest selects function to respend
func handleRequest(context context.Context, key string, writer http.ResponseWriter, request *http.Request,
	funcGet func(context.Context, string, http.ResponseWriter, *http.Request),
	funcPost func(context.Context, string, http.ResponseWriter, *http.Request),
	funcPut func(context.Context, string, http.ResponseWriter, *http.Request),
	funcDelete func(context.Context, string, http.ResponseWriter, *http.Request)) {
	switch request.Method {
	case MethodGet:
		funcGet(context, key, writer, request)
		return
	case MethodPost:
		funcPost(context, key, writer, request)
		return
	case MethodPut:
		funcPut(context, key, writer, request)
		return
	case MethodDelete:
		funcDelete(context, key, writer, request)
		return
	default:
		RespondMethodNotAllowed(context, key, writer, request)
	}
}

// RespondMethodNotAllowed responds as HTTP 405 Method Mot Allowed
func RespondMethodNotAllowed(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
	RespondError(context, writer, request, http.StatusMethodNotAllowed, fmt.Sprintf("%s Method Not Allowed to %s", request.Method, request.URL))
}

// RespondError responds specified error
func RespondError(context context.Context, writer http.ResponseWriter, request *http.Request, status int, args ...interface{}) {
	messages := make([]string, len(args))
	for index, arg := range args {
		messages[index] = fmt.Sprint(arg)
	}
	data := map[string]interface{}{
		"error": map[string][]string{
			"messages": messages,
		},
	}

	buffer := new(bytes.Buffer)
	err := json.NewEncoder(buffer).Encode(data)
	if err != nil {
		log.Errorf(context, "%#v", err)
		RespondError(context, writer, request, http.StatusInternalServerError, err.Error())
	} else {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		http.Error(writer, buffer.String(), status)
	}
}

// Respond responds specified interface as JSON
func Respond(context context.Context, writer http.ResponseWriter, request *http.Request, data interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		log.Errorf(context, "%#v", err)
		RespondError(context, writer, request, http.StatusInternalServerError, err.Error())
	}
}

// RegieterMapping registers mapping of relation between model and URL
func RegieterMapping(prefix string, moduleName string, modelName string,
	funcGet func(context.Context, string, http.ResponseWriter, *http.Request),
	funcPost func(context.Context, string, http.ResponseWriter, *http.Request),
	funcPut func(context.Context, string, http.ResponseWriter, *http.Request),
	funcDelete func(context.Context, string, http.ResponseWriter, *http.Request)) {
	urlMappings = append(urlMappings, urlMapping{
		prefix:     prefix,
		moduleName: moduleName,
		modelName:  modelName,
		funcGet:    funcGet,
		funcPost:   funcPost,
		funcPut:    funcPut,
		funcDelete: funcDelete,
	})
}

// Process processes REST operation
func Process(context context.Context, writer http.ResponseWriter, request *http.Request) {
	mapping, key, err := selectMapping(context, request.URL)
	if err != nil {
		if err.Error() == ErrorURLMappingUnMatch {
			RespondError(context, writer, request, http.StatusNotFound, fmt.Sprintf("%s is not found on this server.", request.URL.Path))
			return
		}
		RespondError(context, writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	handleRequest(context, key, writer, request, mapping.funcGet, mapping.funcPost, mapping.funcPut, mapping.funcDelete)
}
