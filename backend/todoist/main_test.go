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
	defer context.Close()
	if err != nil {
		t.Fatal(err)
	}

	// Try to get new Configuration
	config, err = ReadConfiguration(context)
	if err != nil {
		t.Fatal(err)
	}

	// try to push Configuration
	config.APIKey = apiKey
	config2, err = WriteConfiguration(context, config)
	if config2.APIKey != apiKey {
		t.Errorf("got %v\nwant %v", config2.APIKey, apiKey)
	}

	// try to get new Configuration
	config, err = ReadConfiguration(context)
	if config.APIKey != apiKey {
		t.Errorf("got %v\nwant %v", config2.APIKey, apiKey)
	}
}
