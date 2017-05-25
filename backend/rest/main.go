package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"appengine"
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
	funcGet func(string, http.ResponseWriter, *http.Request, appengine.Context)
	// Function on POST Method
	funcPost func(string, http.ResponseWriter, *http.Request, appengine.Context)
	// Function on PUT Method
	funcPut func(string, http.ResponseWriter, *http.Request, appengine.Context)
	// Function on DELETE Method
	funcDelete func(string, http.ResponseWriter, *http.Request, appengine.Context)
}

// Initialize this package
func init() {
	urlMappings = make([]urlMapping, 0, 10)
}

// selectMapping selects registered urlMapping matches specified URL
func selectMapping(url *url.URL, context appengine.Context) (*urlMapping, string, error) {
	paths := strings.Split(strings.Trim(url.Path, URLPathSeparator), URLPathSeparator)
	for _, mapping := range urlMappings {
		if mapping.prefix == paths[0] && mapping.moduleName == paths[1] && mapping.modelName == paths[2] {
			if len(paths) > 3 {
				return &mapping, paths[3], nil
			}
			return &mapping, "", nil
		}
	}
	context.Infof("mapping of %s not found in %#v", url.Path, urlMappings)
	return new(urlMapping), "", errors.New(ErrorURLMappingUnMatch)
}

// HandleRequest selects function to respend
func handleRequest(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context, funcGet func(string, http.ResponseWriter, *http.Request, appengine.Context), funcPost func(string, http.ResponseWriter, *http.Request, appengine.Context), funcPut func(string, http.ResponseWriter, *http.Request, appengine.Context), funcDelete func(string, http.ResponseWriter, *http.Request, appengine.Context)) {
	switch request.Method {
	case MethodGet:
		funcGet(key, writer, request, context)
		return
	case MethodPost:
		funcPost(key, writer, request, context)
		return
	case MethodPut:
		funcPut(key, writer, request, context)
		return
	case MethodDelete:
		funcDelete(key, writer, request, context)
		return
	default:
		RespondMethodNotAllowed(writer, request, context)
	}
}

// RespondMethodNotAllowed responds as HTTP 405 Method Mot Allowed
func RespondMethodNotAllowed(writer http.ResponseWriter, request *http.Request, context appengine.Context) {
	RespondError(writer, request, context, http.StatusMethodNotAllowed, fmt.Sprintf("%s Method Not Allowed to %s", request.Method, request.URL))
}

// RespondError responds specified error
func RespondError(writer http.ResponseWriter, request *http.Request, context appengine.Context, status int, args ...interface{}) {
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
		context.Errorf("%#v", err)
		RespondError(writer, request, context, http.StatusInternalServerError, err.Error())
	} else {
		writer.Header().Set("Content-Type", "application/json; charset=utf-8")
		http.Error(writer, buffer.String(), status)
	}
}

// Respond responds specified interface as JSON
func Respond(writer http.ResponseWriter, request *http.Request, context appengine.Context, data interface{}) {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	err := json.NewEncoder(writer).Encode(data)
	if err != nil {
		context.Errorf("%#v", err)
		RespondError(writer, request, context, http.StatusInternalServerError, err.Error())
	}
}

// RegieterMapping registers mapping of relation between model and URL
func RegieterMapping(prefix string, moduleName string, modelName string, funcGet func(string, http.ResponseWriter, *http.Request, appengine.Context), funcPost func(string, http.ResponseWriter, *http.Request, appengine.Context), funcPut func(string, http.ResponseWriter, *http.Request, appengine.Context), funcDelete func(string, http.ResponseWriter, *http.Request, appengine.Context)) {
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
func Process(writer http.ResponseWriter, request *http.Request, context appengine.Context) {
	mapping, key, err := selectMapping(request.URL, context)
	if err != nil {
		if err.Error() == ErrorURLMappingUnMatch {
			RespondError(writer, request, context, http.StatusNotFound, fmt.Sprintf("%s is not found on this server.", request.URL.Path))
			return
		}
		RespondError(writer, request, context, http.StatusInternalServerError, err.Error())
		return
	}
	handleRequest(key, writer, request, context, mapping.funcGet, mapping.funcPost, mapping.funcPut, mapping.funcDelete)
}
