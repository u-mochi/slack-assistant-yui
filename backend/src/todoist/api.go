/**
 * todoist/api.go
 * This file implements API of todoist package
 *
 * auther u-mochi
 * license MIT
 */

// Package todoist implements Datamodel of Todoist and client of Todoist API, configurations.
package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"unicode/utf8"

	"rest"

	"appengine"
)

// resopndConfiguration responds Configuration
func resopndConfiguration(writer http.ResponseWriter, request *http.Request, context appengine.Context) {
	configuration, err := GetConfiguration(context)
	if err != nil {
		context.Errorf("Error on get Configuration. %s\n%#v", err.Error(), err)
		rest.RespondError(writer, request, context, http.StatusInternalServerError, err.Error())
	}
	rest.Respond(writer, request, context, []Configuration{configuration})
}

// checkKey checks key, when key is set, return error.
func checkKey(key string) error {
	if utf8.RuneCountInString(key) > 0 {
		return fmt.Errorf("Key is not required. (%s)", key)
	}
	return nil
}

// ReadOperation returns slice of todoist.Configuration
func ReadOperation(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
	if err := checkKey(key); err != nil {
		context.Errorf("Error occurred. (%#v)", err)
		rest.RespondError(writer, request, context, http.StatusBadRequest, err.Error())
		return
	}
	resopndConfiguration(writer, request, context)
}

// UpdateOperation updates specified todoist.Configuration
func UpdateOperation(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
	if err := checkKey(key); err != nil {
		context.Errorf("Error occurred. (%#v)", err)
		rest.RespondError(writer, request, context, http.StatusBadRequest, err.Error())
		return
	}
	var config Configuration
	bufferBody := new(bytes.Buffer)
	bufferBody.ReadFrom(request.Body)
	err := json.Unmarshal(bufferBody.Bytes(), &config)
	if err != nil {
		context.Errorf("Error on decode JSON. %#v: JSON %s", err, bufferBody.String())
		rest.RespondError(writer, request, context, http.StatusBadRequest, err.Error())
		return
	}
	_, err = SetConfiguration(context, config)
	if err != nil {
		context.Errorf("Error on put Configuration. %#v (%#v)", err, config)
		rest.RespondError(writer, request, context, http.StatusInternalServerError, err.Error())
		return
	}
	resopndConfiguration(writer, request, context)
}
