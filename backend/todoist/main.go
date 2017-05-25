// Package todoist implements Datamodel of Todoist and client of Todoist API, configurations.
//
// auther u-mochi
// license MIT
package todoist

import (
	"time"

	"appengine/datastore"

	"appengine"
)

const (
	// KindConfigure means Datastore kind of Configuration struct.
	kindStringConfiguration = "todoist_configure"
	// KeyConfigure means string key of Datastore for Configuration.
	keyStringConfiguration = "key_todoist_configure"
)

var (
	// keyConfiguration is datastore.Key to manage Configuration.
	keyConfiguration *datastore.Key
)

// Configuration means configurations of todoist packeage
type Configuration struct {
	// API key of Todoist API
	APIKey string "json:api_key"
	// Updated date of this struct
	UpdateDate time.Time "json:update_date"
}

// getKeyConfiguration returns datastore.Key to manage Configuration
func getKeyConfiguration(c appengine.Context) *datastore.Key {
	if keyConfiguration == nil {
		keyConfiguration = datastore.NewKey(c, kindStringConfiguration, keyStringConfiguration, 0, nil)
	}
	return keyConfiguration
}

// SetConfiguration writes Configuration to Datastore
func SetConfiguration(c appengine.Context, config Configuration) (Configuration, error) {
	config.UpdateDate = time.Now()
	_, err := datastore.Put(c, getKeyConfiguration(c), &config)
	return config, err
}

// GetConfiguration reads Configuration from Datastore
func GetConfiguration(c appengine.Context) (Configuration, error) {
	config := Configuration{}
	err := datastore.Get(c, getKeyConfiguration(c), &config)
	if err == datastore.ErrNoSuchEntity {
		config, err = SetConfiguration(c, config)
	}
	return config, err
}
