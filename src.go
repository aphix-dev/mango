package mango

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	PUB = iota
	PRIV
	MUT
)

/**
TAGS
access:
	"mut" -> modifiable by client on creation/update
	"public" -> accessible to all clients
	"private" -> accessible to only the owning client
	"serverOnly" -> accessible only to the server
*/

func Make[Type any](in Type, filterType int) Type {
	var a Type
	pub := reflect.New(reflect.TypeOf(a))
	v := reflect.ValueOf(in)
	ValidateStruct(reflect.TypeOf(a))
	switch filterType {
	case PUB:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			if strings.Contains(accessTag, "public") || accessTag == "" {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
				fmt.Printf("%v [=] %v\n", v.Type().Field(i).Name, v.Field(i))
				fmt.Printf("pub = %v\n", reflect.Indirect(pub))
			}
		}
	case PRIV:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			if strings.Contains(accessTag, "private") || strings.Contains(accessTag, "public") || accessTag == "" {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
			}
		}
	case MUT:
		for i := 0; i < v.NumField(); i++ {
			accessTag := v.Type().Field(i).Tag.Get("access")
			if strings.Contains(accessTag, "mut") {
				reflect.Indirect(pub).Field(i).Set(v.Field(i))
			}
		}
	}

	g := pub.Interface()

	return *(g.(*Type))
}

func ValidateStruct(s reflect.Type) {
	for i := 0; i < s.NumField(); i++ {
		accessTag := s.Field(i).Tag.Get("access")
		if (strings.Contains(accessTag, "public") && strings.Contains(accessTag, "private")) ||
			(strings.Contains(accessTag, "serverOnly") && (strings.Contains(accessTag, "public") || strings.Contains(accessTag, "private"))) {
			panic("struct is invalid!")
		}
	}
	fmt.Println("struct is valid")
}

func MakeTruePointer() *bool {
	a := new(bool)
	*a = true
	return a
}

func MakeFalsePointer() *bool {
	a := new(bool)
	*a = false
	return a
}
