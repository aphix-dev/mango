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
TAGS
access:
	"create" -> field is modifiable by client on creation
	"update" -> field is modifiable by client on update
	"pub" -> field is accessible to all clients
	"priv" -> field is accessible to only the owning client

NOTE: fields that are both public and private should be marked with both the "pub" and "priv" tags
*/

/*
`in` => the variable whose fields should be stripped
`filterType` => determines what fields of `in` should be trimmed, based on access tag
*/
func Trim[Type any](rawStruct Type, filterType int) Type {
	var a Type
	pub := reflect.New(reflect.TypeOf(a))
	// set `v` to the value stored in `in`
	v := reflect.ValueOf(rawStruct)
	ValidateStruct(reflect.TypeOf(a))
	switch filterType {
	case PUB:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			// if the field is tagged as "public", return it as a public field
			if strings.Contains(accessTag, "pub") {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
				fmt.Printf("%v [=] %v\n", v.Type().Field(i).Name, v.Field(i))
				fmt.Printf("pub = %v\n", reflect.Indirect(pub))
			}
		}
	case PRIV:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			// if the field is tagged as "private", return it as a private field
			if strings.Contains(accessTag, "priv") {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
			}
		}
	case CREATE:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			// if the field is tagged as "create", return it as a create field
			if strings.Contains(accessTag, "create") {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
			}
		}
	case UPDATE:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			// if the field is tagged as "update", return it as an update field
			if strings.Contains(accessTag, "update") {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
			}
		}
	}

	g := pub.Interface()

	return g.(Type)
}

/*
Checks if the type of `s` is tagged properly
*/
func ValidateStruct(s reflect.Type) {
	for i := 0; i < s.NumField(); i++ {
		accessTag := s.Field(i).Tag.Get("access")
		if (strings.Contains(accessTag, "public") && strings.Contains(accessTag, "private")) ||
			(strings.Contains(accessTag, "serverOnly") && (strings.Contains(accessTag, "public") || strings.Contains(accessTag, "private"))) {
		}
	}
}
