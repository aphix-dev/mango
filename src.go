package mango

import (
	"fmt"
	"reflect"
	"strings"
)

/*
DEFAULT TAGS

access:

	"create" -> field is modifiable by client on creation
	"update" -> field is modifiable by client on update
	"pub" -> field is accessible to all clients
	"priv" -> field is accessible to only the owning client

NOTE: fields that are both public and private should be marked with both the "pub" and "priv" tags
*/
const (
	CREATE = iota - 4
	UPDATE
	PRIV
	PUB
)

// MangoConfig is used to set up different ways to
// trim a struct based on access privileges.
// It can be extended to work with custom access tags.
type MangoConfig struct {
	Filters map[int]string
	Log     bool
}

// DefaultConfig returns the default MangoConfig, which is useful for general REST purposes
func DefaultConfig() MangoConfig {
	return MangoConfig{
		Filters: map[int]string{
			CREATE: "create",
			UPDATE: "update",
			PUB:    "pub",
			PRIV:   "priv",
		},
	}
}

// Extend adds additional filters to an already existing MangoConfig
func (conf MangoConfig) Extend(newFilters map[int]string) MangoConfig {
	for filterType, accessTag := range newFilters {
		conf.Filters[filterType] = accessTag
	}

	return conf
}

// EnableLogs turns on mango logging
func (conf MangoConfig) EnableLogs() MangoConfig {
	conf.Log = true
	return conf
}

// Trim uses a MangoConfig to trim a struct to a form specified by filterType.
// filterType will be used to look for an access tag in the MangoConfig.Filters.
// The returned value is the struct with only fields with the specified `access` tag present.
func Trim[Type any](rawStruct Type, filterType int, conf MangoConfig) (trimmedStruct Type) {
	// create a blank variable of the same type as `rawStruct`
	// to be returned after being filled with the trimmed fields
	trimmed := reflect.New(reflect.TypeOf(rawStruct))

	// create a reflect.Value representation of `rawStruct`
	baseData := reflect.ValueOf(rawStruct)

	// fetch the accessTag string to trim to
	accessTag := conf.Filters[filterType]

	// iterate through every field in `baseData`
	for i := 0; i < baseData.NumField(); i++ {
		curAccessTag := baseData.Type().Field(i).Tag.Get("access")
		if strings.Contains(curAccessTag, accessTag) {
			reflect.Indirect(trimmed).Field(i).Set(baseData.Field(i))

			// log data if logs are set in `conf`
			if conf.Log {
				fmt.Printf("mango: keep: [%v] = %v\n", baseData.Type().Field(i).Name, baseData.Field(i))
			}

		}
	}

	asInterface := trimmed.Interface()

	return *(asInterface.(*Type))
}
