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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"rest"
	"unicode/utf8"

	"google.golang.org/appengine/log"
)

// resopndConfiguration responds Configuration
func resopndConfiguration(context context.Context, writer http.ResponseWriter, request *http.Request) {
	configuration, err := GetConfiguration(context)
	if err != nil {
		log.Errorf(context, "Error on get Configuration. %s\n%#v", err.Error(), err)
		rest.RespondError(context, writer, request, http.StatusInternalServerError, err.Error())
	}
	rest.Respond(context, writer, request, []Configuration{configuration})
}

// checkKey checks key, when key is set, return error.
func checkKey(key string) error {
	if utf8.RuneCountInString(key) > 0 {
		return fmt.Errorf("Key is not required. (%s)", key)
	}
	return nil
}

// ReadOperation returns slice of todoist.Configuration
func ReadOperation(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
	if err := checkKey(key); err != nil {
		log.Errorf(context, "Error occurred. (%#v)", err)
		rest.RespondError(context, writer, request, http.StatusBadRequest, err.Error())
		return
	}
	resopndConfiguration(context, writer, request)
}

// UpdateOperation updates specified todoist.Configuration
func UpdateOperation(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
	if err := checkKey(key); err != nil {
		log.Errorf(context, "Error occurred. (%#v)", err)
		rest.RespondError(context, writer, request, http.StatusBadRequest, err.Error())
		return
	}
	var config Configuration
	bufferBody := new(bytes.Buffer)
	bufferBody.ReadFrom(request.Body)
	err := json.Unmarshal(bufferBody.Bytes(), &config)
	if err != nil {
		log.Errorf(context, "Error on decode JSON. %#v: JSON %s", err, bufferBody.String())
		rest.RespondError(context, writer, request, http.StatusBadRequest, err.Error())
		return
	}
	_, err = SetConfiguration(context, config)
	if err != nil {
		log.Errorf(context, "Error on put Configuration. %#v (%#v)", err, config)
		rest.RespondError(context, writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	resopndConfiguration(context, writer, request)
}
