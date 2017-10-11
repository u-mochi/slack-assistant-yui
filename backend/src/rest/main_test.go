/**
 * rest/main_test.go
 * This file implements REST framework test of rest package
 *
 * auther u-mochi
 * license MIT
 */

// Package rest implements REST framework of Slack Assistant Yui.
package rest

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"context"

	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/log"
)

// context is Context of this test
var (
	// context is Context of this test
	ctx context.Context
	// closure that must be called when the Context is no longer required.
	done func()
	// errorOnInit is error on init()
	errorOnInit error
)

// Initialize this tests
func init() {
	ctx, done, errorOnInit = aetest.NewContext()
	defer done()
}

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

// TestselectMapping tests selectMapping
func TestSelectMapping(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}

	const (
		prefix     = "api"
		moduleName = "testmodule"
		modelName  = "testmodel"
		modelName2 = "testmodel2"
		keyValue   = "somekey"
		path1      = "/api/testmodule/testmodel/"
		path2      = "/api/testmodule2/testmodel/"
	)

	//success test
	urlMappings = nil
	var url1 = url.URL{Path: path1}
	RegieterMapping(prefix, moduleName, modelName,
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {})
	RegieterMapping(prefix, moduleName, modelName2,
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {},
		func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {})
	mapping, key, err := selectMapping(ctx, &url1)
	if err != nil {
		t.Fatal(err)
	}
	assertsEquals(t, prefix, mapping.prefix)
	assertsEquals(t, moduleName, mapping.moduleName)
	assertsEquals(t, modelName, mapping.modelName)
	assertsEquals(t, key, "")

	url1.Path += keyValue + URLPathSeparator
	mapping, key, err = selectMapping(ctx, &url1)
	if err != nil {
		t.Fatal(err)
	}
	assertsEquals(t, prefix, mapping.prefix)
	assertsEquals(t, moduleName, mapping.moduleName)
	assertsEquals(t, modelName, mapping.modelName)
	assertsEquals(t, key, keyValue)

	//fail test
	url1.Path = path2
	mapping, key, err = selectMapping(ctx, &url1)
	if err == nil {
		t.Fatal("Error is not set.")
	}

	log.Infof(ctx, "TestselectMapping passed.")
}

// TestRespondMethodNotAllowed tests rest.RespondMethodNotAllowed
func TestRespondMethodNotAllowed(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}

	request, response, err := createRequest("test", MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	RespondMethodNotAllowed(ctx, "", response, request)
	assertsEquals(t, http.StatusMethodNotAllowed, response.Code)
	log.Debugf(ctx, "Response: %#v\nBody: %s", response, response.Body.String())
	log.Infof(ctx, "TestRespondMethodNotAllowed passed.")
}

// TestRespondMethodNotAllowed tests rest.RespondError
func TestRespondError(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}

	request, response, err := createRequest("test", MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	code := http.StatusBadRequest
	RespondError(ctx, response, request, code, "arg1", errors.New("test"))
	assertsEquals(t, http.StatusBadRequest, code)
	log.Debugf(ctx, "Response: %#v\nBody: %s", response, response.Body.String())
	log.Infof(ctx, "TestRespondError passed.")
}

// TestRespondMethodNotAllowed tests rest.RespondError
func TestRespond(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}

	request, response, err := createRequest("test", MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	type sampleStruct struct {
		Prop1 string
		Prop2 int
	}
	testStuct := sampleStruct{Prop1: "Test", Prop2: 3}

	Respond(ctx, response, request, testStuct)
	log.Debugf(ctx, "Response: %#v\nBody: %s", response, response.Body.String())
	var decodedStruct sampleStruct
	err2 := json.Unmarshal(response.Body.Bytes(), &decodedStruct)
	if err2 != nil {
		t.Fatal(err2)
	}
	assertsEquals(t, testStuct.Prop1, decodedStruct.Prop1)
	assertsEquals(t, testStuct.Prop2, decodedStruct.Prop2)
	log.Infof(ctx, "TestRespond passed.")
}

func TestHandleRequest(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}

	var (
		isFuncGetCalled    bool
		isFuncPostCalled   bool
		isFuncPutCalled    bool
		isFuncDeleteCalled bool
		funcReset          = func() {
			isFuncGetCalled = false
			isFuncPostCalled = false
			isFuncPutCalled = false
			isFuncDeleteCalled = false
		}
		funcGet = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncGetCalled = true
		}
		funcPost = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPostCalled = true
		}
		funcPut = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPutCalled = true
		}
		funcDelete = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncDeleteCalled = true
		}
	)

	request, response, err := createRequest("test", MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}

	// Test Option
	handleRequest(ctx, "", response, request, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, http.StatusMethodNotAllowed, response.Code)
	assertsEquals(t, false, isFuncGetCalled)
	assertsEquals(t, false, isFuncPostCalled)
	assertsEquals(t, false, isFuncPutCalled)
	assertsEquals(t, false, isFuncDeleteCalled)

	// test Get
	funcReset()
	request, response, err = createRequest("test", MethodGet, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	handleRequest(ctx, "", response, request, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncGetCalled)
	assertsEquals(t, false, isFuncPostCalled)
	assertsEquals(t, false, isFuncPutCalled)
	assertsEquals(t, false, isFuncDeleteCalled)

	// test Post
	funcReset()
	request, response, err = createRequest("test", MethodPost, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	handleRequest(ctx, "", response, request, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, false, isFuncGetCalled)
	assertsEquals(t, true, isFuncPostCalled)
	assertsEquals(t, false, isFuncPutCalled)
	assertsEquals(t, false, isFuncDeleteCalled)

	// test Put
	funcReset()
	request, response, err = createRequest("test", MethodPut, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	handleRequest(ctx, "", response, request, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, false, isFuncGetCalled)
	assertsEquals(t, false, isFuncPostCalled)
	assertsEquals(t, true, isFuncPutCalled)
	assertsEquals(t, false, isFuncDeleteCalled)

	// test Delete
	funcReset()
	request, response, err = createRequest("test", MethodDelete, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	handleRequest(ctx, "", response, request, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, false, isFuncGetCalled)
	assertsEquals(t, false, isFuncPostCalled)
	assertsEquals(t, false, isFuncPutCalled)
	assertsEquals(t, true, isFuncDeleteCalled)
}

func TestRegieterMapping(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}
	const (
		prefix     = "prefix"
		moduleName = "moduleName"
		modelName  = "modelName"
	)
	var (
		isFuncGetCalled    bool
		isFuncPostCalled   bool
		isFuncPutCalled    bool
		isFuncDeleteCalled bool
		funcReset          = func() {
			isFuncGetCalled = false
			isFuncPostCalled = false
			isFuncPutCalled = false
			isFuncDeleteCalled = false
		}
		funcGet = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncGetCalled = true
		}
		funcPost = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPostCalled = true
		}
		funcPut = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPutCalled = true
		}
		funcDelete = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncDeleteCalled = true
		}
		createdURLMapping urlMapping
	)

	funcReset()
	request, response, err := createRequest("test", MethodPost, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	urlMappings = nil

	assertsEquals(t, 0, len(urlMappings))
	RegieterMapping(prefix, moduleName, modelName, funcGet, funcPost, funcPut, funcDelete)
	assertsEquals(t, 1, len(urlMappings))
	createdURLMapping = urlMappings[0]
	assertsEquals(t, prefix, createdURLMapping.prefix)
	assertsEquals(t, moduleName, createdURLMapping.moduleName)
	assertsEquals(t, modelName, createdURLMapping.modelName)
	funcReset()
	createdURLMapping.funcGet(ctx, "", response, request)
	assertsEquals(t, true, isFuncGetCalled)
	funcReset()
	createdURLMapping.funcPost(ctx, "", response, request)
	assertsEquals(t, true, isFuncPostCalled)
	funcReset()
	createdURLMapping.funcPut(ctx, "", response, request)
	assertsEquals(t, true, isFuncPutCalled)
	funcReset()
	createdURLMapping.funcDelete(ctx, "", response, request)
	assertsEquals(t, true, isFuncDeleteCalled)

	log.Infof(ctx, "TestResieterMapping passed.")
}

func TestProcess(t *testing.T) {
	if errorOnInit != nil {
		t.Fatal(errorOnInit)
	}
	const (
		prefix      = "prefix"
		moduleName1 = "moduleName1"
		moduleName2 = "moduleName2"
		modelName1  = "modelName1"
		modelName2  = "modelName2"
		validURL    = URLPathSeparator + prefix + URLPathSeparator + moduleName1 + URLPathSeparator + modelName1 + URLPathSeparator
		invalidURL  = URLPathSeparator + prefix + URLPathSeparator + modelName2 + URLPathSeparator + modelName1 + URLPathSeparator
	)
	var (
		isFuncGetCalled    bool
		isFuncPostCalled   bool
		isFuncPutCalled    bool
		isFuncDeleteCalled bool
		funcReset          = func() {
			isFuncGetCalled = false
			isFuncPostCalled = false
			isFuncPutCalled = false
			isFuncDeleteCalled = false
		}
		funcGet = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncGetCalled = true
		}
		funcPost = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPostCalled = true
		}
		funcPut = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncPutCalled = true
		}
		funcDelete = func(context context.Context, key string, writer http.ResponseWriter, request *http.Request) {
			isFuncDeleteCalled = true
		}
	)
	RegieterMapping(prefix, moduleName1, modelName1, funcGet, funcPost, funcPut, funcDelete)
	RegieterMapping(prefix, moduleName1, modelName2, funcGet, funcPost, funcPut, funcDelete)
	funcReset()

	// Fail test 404
	request, response, err := createRequest(invalidURL, MethodPost, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(ctx, response, request)
	assertsEquals(t, http.StatusNotFound, response.Code)

	// Success test GET
	funcReset()
	request, response, err = createRequest(validURL, MethodGet, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(ctx, response, request)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncGetCalled)

	// Success test POST
	funcReset()
	request, response, err = createRequest(validURL, MethodPost, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(ctx, response, request)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncPostCalled)

	// Success test PUT
	funcReset()
	request, response, err = createRequest(validURL, MethodPut, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(ctx, response, request)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncPutCalled)

	// Success test DELETE
	funcReset()
	request, response, err = createRequest(validURL, MethodDelete, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(ctx, response, request)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncDeleteCalled)

	log.Infof(ctx, "TestResieterMapping passed.")
}
