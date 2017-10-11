/**
 * todoist/main_test.go
 * This file implements model test of todoist package
 *
 * auther u-mochi
 * license MIT
 */

// Package todoist implements Datamodel of Todoist and client of Todoist API, configurations.
package todoist

import (
	"context"
	"net/http"
	"testing"

	"google.golang.org/appengine/aetest"
)

// Test Configure
func TestConfigure(t *testing.T) {
	var (
		context context.Context
		done    func()
		config  Configuration
		config2 Configuration
		err     error
		apiKey  = "12345"
	)

	// Try to get aetest.Context
	context, done, err = aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// Try to get new Configuration
	config, err = GetConfiguration(context)
	if err != nil {
		t.Fatal(err)
	}

	// try to push Configuration
	config.APIKey = apiKey
	config2, err = SetConfiguration(context, config)
	if config2.APIKey != apiKey {
		t.Errorf("got %v\nwant %v", config2.APIKey, apiKey)
	}

	// try to get new Configuration
	config, err = GetConfiguration(context)
	if config.APIKey != apiKey {
		t.Errorf("got %v\nwant %v", config2.APIKey, apiKey)
	}
}

func TestSync(t *testing.T) {
	var (
		response *http.Response
		context  context.Context
		done     func()
		err      error
	)

	// Try to get aetest.Context
	context, done, err = aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	// Try to sync Todoist
	response, err = Sync(context)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Resp: %#v", response)
}
