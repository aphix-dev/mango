package mango

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	CREATE = iota
	UPDATE
	PRIV
	PUB
)

/**
DEFAULT TAGS
access:
	"create" -> field is modifiable by client on creation
	"update" -> field is modifiable by client on update
	"pub" -> field is accessible to all clients
	"priv" -> field is accessible to only the owning client

NOTE: fields that are both public and private should be marked with both the "pub" and "priv" tags
*/

// MangoConfig is used to set up different ways to
// trim a struct based on access privileges.
// It can be extended to work with custom access tags.
type MangoConfig struct {
	Filters map[int]string
	Log     bool
}

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

func (conf MangoConfig) Extend(newFilters map[int]string) MangoConfig {
	for filterType, accessTag := range newFilters {
		conf.Filters[filterType] = accessTag
	}

	return conf
}

func (conf MangoConfig) EnableLogs() MangoConfig {
	conf.Log = true
	return conf
}

func Trim[Type any](rawStruct Type, filterType int, conf MangoConfig) Type {
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
				fmt.Printf("%v [=] %v\n", baseData.Type().Field(i).Name, baseData.Field(i))
				fmt.Printf("baseData = %v\n", reflect.Indirect(baseData))
			}

		}
	}

	asInterface := trimmed.Interface()

	return *(asInterface.(*Type))
}
