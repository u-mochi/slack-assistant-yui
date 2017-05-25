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

	"appengine"

	"appengine/aetest"
)

// context is Context of this test
var (
	// context is Context of this test
	context aetest.Context
	// errorOnInit is error on init()
	errorOnInit error
)

// Initialize this tests
func init() {
	context, errorOnInit = aetest.NewContext(nil)
	defer context.Close()
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
	RegieterMapping(prefix, moduleName, modelName, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {})
	RegieterMapping(prefix, moduleName, modelName2, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {}, func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {})
	mapping, key, err := selectMapping(&url1, context)
	if err != nil {
		t.Fatal(err)
	}
	assertsEquals(t, prefix, mapping.prefix)
	assertsEquals(t, moduleName, mapping.moduleName)
	assertsEquals(t, modelName, mapping.modelName)
	assertsEquals(t, key, "")

	url1.Path += keyValue + URLPathSeparator
	mapping, key, err = selectMapping(&url1, context)
	if err != nil {
		t.Fatal(err)
	}
	assertsEquals(t, prefix, mapping.prefix)
	assertsEquals(t, moduleName, mapping.moduleName)
	assertsEquals(t, modelName, mapping.modelName)
	assertsEquals(t, key, keyValue)

	//fail test
	url1.Path = path2
	mapping, key, err = selectMapping(&url1, context)
	if err == nil {
		t.Fatal("Error is not set.")
	}

	context.Infof("TestselectMapping passed.")
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
	RespondMethodNotAllowed(response, request, context)
	assertsEquals(t, http.StatusMethodNotAllowed, response.Code)
	context.Debugf("Response: %#v\nBody: %s", response, response.Body.String())
	context.Infof("TestRespondMethodNotAllowed passed.")
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
	RespondError(response, request, context, code, "arg1", errors.New("test"))
	assertsEquals(t, http.StatusBadRequest, code)
	context.Debugf("Response: %#v\nBody: %s", response, response.Body.String())
	context.Infof("TestRespondError passed.")
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

	Respond(response, request, context, testStuct)
	context.Debugf("Response: %#v\nBody: %s", response, response.Body.String())
	var decodedStruct sampleStruct
	err2 := json.Unmarshal(response.Body.Bytes(), &decodedStruct)
	if err2 != nil {
		t.Fatal(err2)
	}
	assertsEquals(t, testStuct.Prop1, decodedStruct.Prop1)
	assertsEquals(t, testStuct.Prop2, decodedStruct.Prop2)
	context.Infof("TestRespond passed.")
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
		funcGet = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncGetCalled = true
		}
		funcPost = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPostCalled = true
		}
		funcPut = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPutCalled = true
		}
		funcDelete = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncDeleteCalled = true
		}
	)

	request, response, err := createRequest("test", MethodOption, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}

	// Test Option
	handleRequest("", response, request, context, funcGet, funcPost, funcPut, funcDelete)
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
	handleRequest("", response, request, context, funcGet, funcPost, funcPut, funcDelete)
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
	handleRequest("", response, request, context, funcGet, funcPost, funcPut, funcDelete)
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
	handleRequest("", response, request, context, funcGet, funcPost, funcPut, funcDelete)
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
	handleRequest("", response, request, context, funcGet, funcPost, funcPut, funcDelete)
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
		funcGet = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncGetCalled = true
		}
		funcPost = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPostCalled = true
		}
		funcPut = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPutCalled = true
		}
		funcDelete = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
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
	createdURLMapping.funcGet("", response, request, context)
	assertsEquals(t, true, isFuncGetCalled)
	funcReset()
	createdURLMapping.funcPost("", response, request, context)
	assertsEquals(t, true, isFuncPostCalled)
	funcReset()
	createdURLMapping.funcPut("", response, request, context)
	assertsEquals(t, true, isFuncPutCalled)
	funcReset()
	createdURLMapping.funcDelete("", response, request, context)
	assertsEquals(t, true, isFuncDeleteCalled)

	context.Infof("TestResieterMapping passed.")
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
		funcGet = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncGetCalled = true
		}
		funcPost = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPostCalled = true
		}
		funcPut = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
			isFuncPutCalled = true
		}
		funcDelete = func(key string, writer http.ResponseWriter, request *http.Request, context appengine.Context) {
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
	Process(response, request, context)
	assertsEquals(t, http.StatusNotFound, response.Code)

	// Success test GET
	funcReset()
	request, response, err = createRequest(validURL, MethodGet, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(response, request, context)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncGetCalled)

	// Success test POST
	funcReset()
	request, response, err = createRequest(validURL, MethodPost, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(response, request, context)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncPostCalled)

	// Success test PUT
	funcReset()
	request, response, err = createRequest(validURL, MethodPut, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(response, request, context)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncPutCalled)

	// Success test DELETE
	funcReset()
	request, response, err = createRequest(validURL, MethodDelete, new(strings.Reader))
	if err != nil {
		t.Fatal(errorOnInit)
	}
	Process(response, request, context)
	assertsEquals(t, http.StatusOK, response.Code)
	assertsEquals(t, true, isFuncDeleteCalled)

	context.Infof("TestResieterMapping passed.")
}
