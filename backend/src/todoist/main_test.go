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
	"testing"

	"appengine/aetest"
)

// Test Configure
func TestConfigure(t *testing.T) {
	var (
		context aetest.Context
		config  Configuration
		config2 Configuration
		err     error
		apiKey  = "12345"
	)

	// Try to get aetest.Context
	context, err = aetest.NewContext(nil)
	if err != nil {
		t.Fatal(err)
	}
	defer context.Close()

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
