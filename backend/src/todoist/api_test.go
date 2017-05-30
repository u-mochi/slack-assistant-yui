/**
 * todoist/api_test.go
 * This file implements api test of todoist package
 *
 * auther u-mochi
 * license MIT
 */

// Package todoist implements Datamodel of Todoist and client of Todoist API, configurations.
package todoist

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"rest"
	"strings"
	"testing"

	"appengine/aetest"
)

// createRequest creates Request and ResponseRecorder
func createRequest(url string, method string, body io.Reader) (*http.Request, *httptest.ResponseRecorder, error) {
	request, err := http.NewRequest(method, url, body)
	response := httptest.NewRecorder()
	return request, response, err
}

// assertEquals asserts specified value is same or not
func assertsEquals(t *testing.T, expected interface{}, actual interface{}) {
	if expected != actual {
		t.Errorf("accept: %v\nexpect: %v", actual, expected)
	}
}

// Test ReadOperation
func TestReadOperation(t *testing.T) {
	var (
		request  *http.Request
		response *httptest.ResponseRecorder
		config1  Configuration
		config2  Configuration
		configs  []Configuration
	)

	// Try to get aetest.Context
	context, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer context.Close()

	request, response, err = createRequest("test", rest.MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(err)
	}

	// Try to get new Configuration
	ReadOperation("", response, request, context)
	assertsEquals(t, http.StatusOK, response.Code)
	err = json.Unmarshal(response.Body.Bytes(), &configs)
	if err != nil {
		t.Fatal(err)
		return
	}
	assertsEquals(t, 1, len(configs))
	config1 = configs[0]
	config2, err = GetConfiguration(context)
	if err != nil {
		t.Fatal(err)
		return
	}
	assertsEquals(t, config2.APIKey, config1.APIKey)
	assertsEquals(t, config2.UpdateDate, config1.UpdateDate)

	//Fail test
	response = httptest.NewRecorder()
	ReadOperation("testkey", response, request, context)
	assertsEquals(t, http.StatusBadRequest, response.Code)
}

// Test UpdateOperation
func TestUpdateOperation(t *testing.T) {
	const newAPIKey = "something new api key"
	var (
		request  *http.Request
		response *httptest.ResponseRecorder
		config1  Configuration
		config2  Configuration
		configs  []Configuration
	)

	// Try to get aetest.Context
	context, err := aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer context.Close()

	// Try to get new Configuration
	config1, err = GetConfiguration(context)
	if err != nil {
		t.Fatal(err)
		return
	}
	config1.APIKey = newAPIKey
	bytes, err1 := json.Marshal(config1)
	if err1 != nil {
		t.Fatal(err)
	}
	request, response, err = createRequest("test", rest.MethodOption, strings.NewReader(string(bytes)))
	if err != nil {
		t.Fatal(err)
	}
	UpdateOperation("", response, request, context)

	err = json.Unmarshal(response.Body.Bytes(), &configs)
	if err != nil {
		t.Fatal(err)
		return
	}
	assertsEquals(t, 1, len(configs))
	config2 = configs[0]
	config1, err = GetConfiguration(context)
	if err != nil {
		t.Fatal(err)
		return
	}
	assertsEquals(t, config2.APIKey, config1.APIKey)
	assertsEquals(t, config2.UpdateDate, config1.UpdateDate)

	//Fail test
	response = httptest.NewRecorder()
	UpdateOperation("testkey", response, request, context)
	assertsEquals(t, http.StatusBadRequest, response.Code)
}
